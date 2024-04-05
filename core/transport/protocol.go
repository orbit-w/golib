package transport

import (
	"context"
	"github.com/orbit-w/golib/core/network"
	"net"
)

/*
   @Author: orbit-w
   @File: protocol
   @2024 4月 周五 18:30
*/

type ConnFactory func(ctx context.Context, _conn net.Conn, maxIncomingPacket uint32, head, body []byte) IServerConn

var factories = map[network.Protocol]ConnFactory{}

func RegProtocol(protocol network.Protocol, factory ConnFactory) {
	if _, ok := factories[protocol]; ok {
		panic("protocol already registered")
	}
	factories[protocol] = factory
}

func dispatchProtocol(protocol network.Protocol) ConnFactory {
	if factory, ok := factories[protocol]; ok {
		return factory
	}
	panic("protocol not registered")
}
