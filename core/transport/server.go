package transport

import (
	"github.com/orbit-w/golib/core/network"
	"time"
)

/*
   @Author: orbit-w
   @File: server
   @2023 11月 周五 17:04
*/

type AcceptorOptions struct {
	MaxIncomingPacket uint32
	IsGzip            bool
}

type IServer interface {
	Stop() error
}

type IServerConn interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
}

func Serve(pStr, host string,
	_handle func(conn IServerConn)) (IServer, error) {
	config := DefaultServerConfig()
	op := config.ToAcceptorOptions()
	protocol := parseProtocol(pStr)
	factory := GetFactory(protocol)
	server := factory()
	if err := server.Serve(host, _handle, op); err != nil {
		return nil, err
	}

	return server, nil
}

type Config struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
}

func (c Config) ToAcceptorOptions() network.AcceptorOptions {
	return network.AcceptorOptions{
		MaxIncomingPacket: c.MaxIncomingPacket,
		IsGzip:            c.IsGzip,
	}
}

func DefaultServerConfig() Config {
	return Config{
		MaxIncomingPacket: network.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       network.ReadTimeout,
		WriteTimeout:      WriteTimeout,
	}
}

func parseProtocol(p string) network.Protocol {
	switch p {
	case "tcp":
		return network.TCP
	case "udp":
		return network.UDP
	case "kcp":
		return network.KCP
	default:
		return network.TCP
	}
}
