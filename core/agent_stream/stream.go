package agent_stream

import (
	"fmt"
	"github.com/orbit-w/golib/bases/packet"
	"github.com/orbit-w/golib/core/network"
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
	Recv() (data []byte, err error)
	Close() error
}

type AgentStream struct {
	conn         net.Conn
	codec        *network.Codec
	r            *network.BlockReceiver
	handleStream func(stream IStream) error
}

func NewAgentStream(handle func(stream IStream) error, maxIncomingPacket uint32, isGzip bool, readTimeout time.Duration) *AgentStream {
	return &AgentStream{
		codec:        network.NewCodec(maxIncomingPacket, isGzip, readTimeout),
		r:            network.NewBlockReceiver(),
		handleStream: handle,
	}
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
		stream.r.OnClose(network.ErrCanceled)
		if err != nil {
			if err == io.EOF || network.IsClosedConnError(err) {
				//连接正常断开
			} else {
				log.Println(fmt.Errorf("[AgentStream] stream disconnected: %s", err.Error()))
			}
		}
	}()

	for {
		in, err = stream.codec.BlockDecodeBody(conn, head, body)
		if err != nil {
			if err == io.EOF {
				log.Println("connection is closed")
			} else {
				log.Printf("read error: %v", err)
			}
			return
		}

		stream.r.Put(in.Remain(), nil)
	}
}
