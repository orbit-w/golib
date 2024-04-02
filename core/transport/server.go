package transport

import (
	"github.com/orbit-w/golib/core/network"
	"net"
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

func Serve(p string, listener net.Listener,
	_handle func(conn IServerConn) error, ops ...network.AcceptorOptions) IServer {
	
	server := new(network.Server)
	protocol := parseProtocol(p)
	server.Serve(protocol, listener, func(conn network.IServerConn) error {
		return _handle(conn)
	}, ops...)
	return server
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
