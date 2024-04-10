package agent_stream

import (
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	network2 "github.com/orbit-w/golib/modules/net/network"
	"io"
	"log"
	"net"
	"runtime/debug"
	"time"
)

/*
   @Author: orbit-w
   @File: stream
   @2024 4月 周日 11:17
*/

type IStream interface {
	Send(body []byte) error
	SendPack(body packet.IPacket) error
	Recv() (data []byte, err error)
	Close() error
}

type AgentStream struct {
	conn         net.Conn
	codec        *network2.Codec
	r            *network2.BlockReceiver
	writeTimeout time.Duration
	readTimeout  time.Duration
}

func NewAgentStream(_conn net.Conn, maxIncoming uint32, isGzip bool, wt, rt time.Duration) *AgentStream {
	return &AgentStream{
		conn:         _conn,
		codec:        network2.NewCodec(maxIncoming, isGzip, rt),
		r:            network2.NewBlockReceiver(),
		writeTimeout: wt,
	}
}

func (stream *AgentStream) Send(body []byte) error {
	out, err := stream.codec.EncodeBodyRaw(body)
	if err != nil {
		return err
	}
	defer out.Return()
	if err = stream.conn.SetWriteDeadline(time.Now().Add(stream.writeTimeout)); err != nil {
		return err
	}
	_, err = stream.conn.Write(out.Data())
	return err
}

func (stream *AgentStream) SendPack(body packet.IPacket) error {
	out, err := stream.codec.EncodeBody(body)
	if err != nil {
		return err
	}
	defer out.Return()
	if err = stream.conn.SetWriteDeadline(time.Now().Add(stream.writeTimeout)); err != nil {
		return err
	}
	_, err = stream.conn.Write(out.Data())
	return err
}

func (stream *AgentStream) Recv() (data []byte, err error) {
	return stream.r.Recv()
}

func (stream *AgentStream) Close() error {
	if stream.conn != nil {
		_ = stream.conn.Close()
	}
	return nil
}

func (stream *AgentStream) handleLoop(conn net.Conn, head, body []byte) {
	var (
		err error
		in  packet.IPacket
	)

	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println("stack: ", string(debug.Stack()))
		}

		if conn != nil {
			_ = conn.Close()
		}
		if err != nil {
			if err == io.EOF || network2.IsClosedConnError(err) {
				//连接正常断开
				stream.r.OnClose(network2.ErrCanceled)
			} else {
				sErr := fmt.Errorf("[AgentStream] stream disconnected error: %s ", err.Error())
				stream.r.OnClose(sErr)
				log.Println(sErr)
			}
		} else {
			stream.r.OnClose(network2.ErrCanceled)
		}
	}()

	for {
		in, err = stream.codec.BlockDecodeBody(conn, head, body)
		if err != nil {
			break
		}

		stream.r.Put(in.Remain(), nil)
	}
}
