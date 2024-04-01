package transport

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
	Write(data []byte) error
	Recv() ([]byte, error)
	Close() error
}

type IServerConn interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
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
