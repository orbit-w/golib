package transport

import (
	"context"
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/core/network"
	"github.com/orbit-w/golib/modules/wrappers/sender_wrapper"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

/*
   @Author: orbit-w
   @File: tcp_server
   @2023 11月 周日 21:03
*/

type TcpServerConn struct {
	authed bool
	conn   net.Conn
	codec  *network.Codec
	ctx    context.Context
	cancel context.CancelFunc
	sw     *sender_wrapper.SenderWrapper
	buf    *ControlBuffer
	r      *network.BlockReceiver
}

func NewTcpServerConn(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte) IServerConn {
	if ctx == nil {
		ctx = context.Background()
	}
	cCtx, cancel := context.WithCancel(ctx)
	ts := &TcpServerConn{
		conn:   _conn,
		codec:  network.NewCodec(maxIncomingPacket, false, ReadTimeout),
		ctx:    cCtx,
		cancel: cancel,
		r:      network.NewBlockReceiver(),
	}

	sw := sender_wrapper.NewSender(ts.SendData)
	ts.sw = sw
	ts.buf = NewControlBuffer(maxIncomingPacket, ts.sw)

	go ts.HandleLoop(head, body)
	return ts
}

func (ts *TcpServerConn) Send(data []byte) (err error) {
	pack := packHeadByte(data, TypeMessageRaw)
	err = ts.buf.Set(pack)
	pack.Return()
	return
}

func (ts *TcpServerConn) Recv() ([]byte, error) {
	return ts.r.Recv()
}

func (ts *TcpServerConn) Close() error {
	return ts.conn.Close()
}

// SendData implicitly call body.Return
// coding: size<int32> | gzipped<bool> | body<bytes>
func (ts *TcpServerConn) SendData(body packet.IPacket) error {
	pack, err := ts.codec.EncodeBody(body, false)
	if err != nil {
		return err
	}
	defer pack.Return()
	if err = ts.conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		return err
	}
	_, err = ts.conn.Write(pack.Data())
	return err
}

func (ts *TcpServerConn) HandleLoop(header, body []byte) {
	var (
		err  error
		data packet.IPacket
	)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("stack: ", string(debug.Stack()))
		}
		ts.r.OnClose(ErrCanceled)
		ts.buf.OnClose()
		if ts.conn != nil {
			_ = ts.conn.Close()
		}
		if err != nil {
			if err == io.EOF || IsClosedConnError(err) {
				//连接正常断开
			} else {
				log.Println(fmt.Errorf("[TcpServerConn] tcp_conn disconnected: %s", err.Error()))
			}
		}
	}()

	for {
		data, err = ts.codec.BlockDecodeBody(ts.conn, header, body)
		if err != nil {
			return
		}
		if err = ts.OnData(data); err != nil {
			//TODO: 错误处理？
			return
		}
	}
}

func (ts *TcpServerConn) OnData(data packet.IPacket) error {
	defer data.Return()
	for len(data.Remain()) > 0 {
		if bytes, err := data.ReadBytes32(); err == nil {
			reader := packet.Reader(bytes)
			_ = ts.HandleData(reader)
		}
	}
	return nil
}

func (ts *TcpServerConn) HandleData(in packet.IPacket) error {
	err := unpackHeadByte(in, func(head int8, data []byte) {
		switch head {
		case TypeMessageHeartbeat:
			ack := packHeadByte(nil, TypeMessageHeartbeatAck)
			_ = ts.buf.Set(ack)
			ack.Return()
		case TypeMessageHeartbeatAck:
		default:
			ts.r.Put(data, nil)
		}
	})
	return err
}
