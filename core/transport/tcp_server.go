package transport

import (
	"context"
	"github.com/orbit-w/golib/core/network"
	"net"
)

/*
   @Author: orbit-w
   @File: tcp_server
   @2024 4月 周二 16:39
*/

func init() {
	RegisterFactory(network.TCP, func() ITransportServer {
		return &TcpServer{}
	})
}

type TcpServer struct {
	server *network.Server
}

func (t *TcpServer) Serve(host string, _handle func(conn IServerConn), op network.AcceptorOptions) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	server := new(network.Server)
	server.Serve(network.TCP, listener, func(ctx context.Context, generic net.Conn, maxIncomingPacket uint32, head, body []byte) {
		conn := NewTcpServerConn(ctx, generic, maxIncomingPacket, head, body)
		defer func() {
			_ = conn.Close()
		}()
		_handle(conn)
	}, op)
	t.server = server
	return nil
}

// Stop stops the server
// 具有可重入性且线程安全, 这意味着这个方法可以被并发多次调用，而不会影响程序的状态或者产生不可预期的结果
func (t *TcpServer) Stop() error {
	if t.server != nil {
		_ = t.server.Stop()
	}
	return nil
}
