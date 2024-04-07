package agent_stream

import (
	"github.com/orbit-w/golib/core/network"
	"github.com/smallnest/rpcx/server"
	"time"
)

/*
   @Author: orbit-w
   @File: serve
   @2024 4月 周日 11:20
*/

func Serve(addr string, handle func(stream IStream) error) error {
	conf := DefaultConfig()
	return ServeByConfig(addr, handle, conf)
}

func ServeByConfig(addr string, handle func(stream IStream) error, conf *Config) error {
	s := server.NewServer()
	if err := s.RegisterName("StreamService",
		NewStreamService(handle, conf), ""); err != nil {
		return New(ErrHeadServe, err.Error())
	}
	err := s.Serve("tcp", addr)
	if err != nil {
		return New(ErrHeadServe, err.Error())
	}
	return nil
}

type Config struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxIncomingPacket: network.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       time.Second * 60,
	}
}
