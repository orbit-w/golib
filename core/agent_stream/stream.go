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
	Send(body []byte) error
	SendPack(body packet.IPacket) error
	Recv() (data []byte, err error)
	Close() error
}

type AgentStream struct {
	conn         net.Conn
	conf         *Config
	codec        *network.Codec
	r            *network.BlockReceiver
	writeTimeout time.Duration
}

func NewAgentStream(_conn net.Conn, _conf *Config) *AgentStream {
	return &AgentStream{
		conn:         _conn,
		conf:         _conf,
		codec:        network.NewCodec(_conf.MaxIncomingPacket, _conf.IsGzip, _conf.ReadTimeout),
		r:            network.NewBlockReceiver(),
		writeTimeout: _conf.WriteTimeout,
	}
}

func (stream *AgentStream) Send(body []byte) error {
	conf := stream.conf
	out, err := stream.codec.EncodeBodyRaw(body, conf.IsGzip)
	if err != nil {
		return err
	}
	defer out.Return()
	if err = stream.conn.SetWriteDeadline(time.Now().Add(stream.conf.WriteTimeout)); err != nil {
		return err
	}
	_, err = stream.conn.Write(out.Data())
	return err
}

func (stream *AgentStream) SendPack(body packet.IPacket) error {
	conf := stream.conf
	out, err := stream.codec.EncodeBody(body, conf.IsGzip)
	if err != nil {
		return err
	}
	defer out.Return()
	if err = stream.conn.SetWriteDeadline(time.Now().Add(stream.conf.WriteTimeout)); err != nil {
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
		stream.r.OnClose(network.ErrCanceled)
		if err != nil {
			if err == io.EOF || network.IsClosedConnError(err) {
				//连接正常断开
			} else {
				log.Println(fmt.Errorf("[AgentStream] stream disconnected: %server", err.Error()))
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
