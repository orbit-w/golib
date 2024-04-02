package network

import (
	"context"
	"net"
)

/*
   @Author: orbit-w
   @File: protocol
   @2024 4月 周二 11:59
*/

type Protocol string

const (
	TCP Protocol = "tcp"
	KCP Protocol = "kcp"
	UDP Protocol = "udp"
)

type ConnFactory func(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte) IServerConn

var factories = map[Protocol]ConnFactory{}

func RegProtocol(protocol Protocol, factory ConnFactory) {
	if _, ok := factories[protocol]; ok {
		panic("protocol already registered")
	}
	factories[protocol] = factory
}

func DispatchProtocol(protocol Protocol) ConnFactory {
	if factory, ok := factories[protocol]; ok {
		return factory
	}
	panic("protocol not registered")
}
