package agent_stream

import (
	"github.com/orbit-w/golib/modules/net/network"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/share"
	"net"
	"time"
)

/*
   @Author: orbit-w
   @File: serve
   @2024 4月 周日 11:20
*/

type Server struct {
	server *server.Server
}

// Serve 以默认配置启动AgentStream服务
func (s *Server) Serve(addr string, handle func(stream IStream) error) error {
	conf := DefaultConfig()
	return s.ServeByConfig(addr, handle, conf)
}

// ServeByConfig 以自定义配置启动AgentStream服务
func (s *Server) ServeByConfig(addr string, handle func(stream IStream) error, conf *Config) error {
	parseConfig(conf)
	s.server = server.NewServer()
	p := server.NewStreamService(addr, func(conn net.Conn, args *share.StreamServiceArgs) {
		_ = NewStreamService(handle, conf).Stream(conn, nil, nil)
	}, nil, 1000)
	s.server.EnableStreamService(share.StreamServiceName, p)
	return nil
}

func (s *Server) Close() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

type Config struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DialTimeout       time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		MaxIncomingPacket: network.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		DialTimeout:       DialTimeout,
	}
}

func parseConfig(conf *Config) {
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = WriteTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = ReadTimeout
	}
	if conf.MaxIncomingPacket == 0 {
		conf.MaxIncomingPacket = network.MaxIncomingPacket
	}
	if conf.DialTimeout == 0 {
		conf.DialTimeout = DialTimeout
	}
}
