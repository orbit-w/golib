package agent_stream

import (
	"github.com/orbit-w/golib/core/network"
	"net"
	"sync"
)

/*
   @Author: orbit-w
   @File: stream_service
   @2024 4月 周一 00:27
*/

type StreamService struct {
	conf     *Config
	handle   func(stream IStream) error
	headPool *sync.Pool
	bodyPool *sync.Pool
}

func NewStreamService(handle func(stream IStream) error, conf *Config) *StreamService {
	return &StreamService{
		conf:     conf,
		handle:   handle,
		headPool: network.NewBufferPool(network.HeadLen),
		bodyPool: network.NewBufferPool(conf.MaxIncomingPacket),
	}
}

func (s *StreamService) Stream(conn net.Conn, _ *string, _ *string) error {
	conf := s.conf
	stream := NewAgentStream(conn, conf.MaxIncomingPacket, conf.IsGzip, conf.WriteTimeout, conf.ReadTimeout)
	headBuf := s.headPool.Get().(*network.Buffer)
	bodyBuf := s.bodyPool.Get().(*network.Buffer)

	defer func() {
		_ = stream.Close()
		s.headPool.Put(headBuf)
		s.bodyPool.Put(bodyBuf)
	}()

	go stream.handleLoop(conn, headBuf.Bytes, bodyBuf.Bytes)
	return s.handle(stream)
}
