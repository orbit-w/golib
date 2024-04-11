package transport

import (
	"github.com/orbit-w/golib/modules/net/network"
)

/*
   @Author: orbit-w
   @File: transport
   @2023 11月 周日 17:01
*/

// IConn represents a virtual connection to a conceptual endpoint
//
// A ClientConn have one actual connections to the endpoint
// based on configuration
type IConn interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
}

type ITransportServer interface {
	Serve(host string, _handle func(conn IConn), op network.AcceptorOptions) error
	Stop() error
}

type DialOption struct {
	RemoteNodeId      string
	CurrentNodeId     string
	MaxIncomingPacket uint32
	IsBlock           bool
	IsGzip            bool
	DisconnectHandler func(nodeId string)
}

type ConnOption struct {
	MaxIncomingPacket uint32
}
