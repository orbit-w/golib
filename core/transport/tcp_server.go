package transport

import (
	"context"
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
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

type TcpServer struct {
	authed bool
	conn   net.Conn
	codec  *NetCodec
	ctx    context.Context
	cancel context.CancelFunc
	sw     *sender_wrapper.SenderWrapper
	buf    *ControlBuffer
	r      *receiver
}

func NewServerConn(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte) IServerConn {
	if ctx == nil {
		ctx = context.Background()
	}
	cCtx, cancel := context.WithCancel(ctx)
	ts := &TcpServer{
		conn:   _conn,
		codec:  NewTcpCodec(maxIncomingPacket, false),
		ctx:    cCtx,
		cancel: cancel,
		r:      newReceiver(),
	}

	sw := sender_wrapper.NewSender(ts.SendData)
	ts.sw = sw
	ts.buf = NewControlBuffer(maxIncomingPacket, ts.sw)

	go ts.HandleLoop(head, body)
	return ts
}

func (ts *TcpServer) Send(data []byte) (err error) {
	pack := packHeadByte(data, TypeMessageRaw)
	err = ts.buf.Set(pack)
	pack.Return()
	return
}

func (ts *TcpServer) Recv() ([]byte, error) {
	return ts.r.read()
}

func (ts *TcpServer) Close() error {
	return ts.conn.Close()
}

// SendData implicitly call body.Return
// coding: size<int32> | gzipped<bool> | body<bytes>
func (ts *TcpServer) SendData(body packet.IPacket) error {
	pack := ts.codec.EncodeBody(body)
	defer pack.Return()
	if err := ts.conn.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		return err
	}
	_, err := ts.conn.Write(pack.Data())
	return err
}

func (ts *TcpServer) HandleLoop(header, body []byte) {
	var (
		err  error
		data packet.IPacket
	)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("stack: ", string(debug.Stack()))
		}
		ts.r.onClose(ErrCanceled)
		ts.buf.OnClose()
		if ts.conn != nil {
			_ = ts.conn.Close()
		}
		if err != nil {
			if err == io.EOF || IsClosedConnError(err) {
				//连接正常断开
			} else {
				log.Println(fmt.Errorf("[TcpServer] tcp_conn disconnected: %s", err.Error()))
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

func (ts *TcpServer) OnData(data packet.IPacket) error {
	defer data.Return()
	for len(data.Remain()) > 0 {
		if bytes, err := data.ReadBytes32(); err == nil {
			reader := packet.Reader(bytes)
			_ = ts.HandleData(reader)
		}
	}
	return nil
}

func (ts *TcpServer) HandleData(in packet.IPacket) error {
	err := unpackHeadByte(in, func(head int8, data []byte) {
		switch head {
		case TypeMessageHeartbeat:
			ack := packHeadByte(nil, TypeMessageHeartbeatAck)
			_ = ts.buf.Set(ack)
			ack.Return()
		case TypeMessageHeartbeatAck:
		default:
			ts.r.put(data, nil)
		}
	})
	return err
}
