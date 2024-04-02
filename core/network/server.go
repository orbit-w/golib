package network

import (
	"context"
	"log"
	"net"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: server
   @2023 11月 周五 17:04
*/

type IServerConn interface {
	Send(data []byte) error
	Recv() ([]byte, error)
	Close() error
}

type Server struct {
	isGzip            bool
	ccu               int32
	maxIncomingPacket uint32
	state             atomic.Uint32
	host              string
	protocol          Protocol
	listener          net.Listener
	rw                sync.RWMutex
	ctx               context.Context
	cancel            context.CancelFunc
	factory           ConnFactory
	bodyPool          *sync.Pool
	headPool          *sync.Pool

	handle func(conn IServerConn) error
}

type AcceptorOptions struct {
	MaxIncomingPacket uint32
	IsGzip            bool
}

func (ins *Server) Serve(p Protocol, listener net.Listener, _handle func(conn IServerConn) error, ops ...AcceptorOptions) {
	op := parseAndWrapOP(ops...)
	ctx, cancel := context.WithCancel(context.Background())
	ins.rw = sync.RWMutex{}
	ins.state.Store(TypeWorking)
	ins.host = ""
	ins.maxIncomingPacket = op.MaxIncomingPacket
	ins.isGzip = op.IsGzip
	ins.ctx = ctx
	ins.cancel = cancel
	ins.handle = _handle
	ins.listener = listener

	ins.protocol = p
	ins.factory = DispatchProtocol(p)

	ins.headPool = NewBufferPool(HeadLen)
	ins.bodyPool = NewBufferPool(ins.maxIncomingPacket)

	go ins.acceptLoop()
}

func (ins *Server) Stop() error {
	if ins.state.CompareAndSwap(TypeWorking, TypeStopped) {
		if ins.cancel != nil {
			ins.cancel()
		}
		if ins.listener != nil {
			_ = ins.listener.Close()
		}
	}
	return nil
}

func (ins *Server) acceptLoop() {
	for {
		conn, err := ins.listener.Accept()
		if err != nil {
			select {
			case <-ins.ctx.Done():
				return
			default:
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}

		ins.handleConn(conn)
	}
}

func (ins *Server) handleConn(conn net.Conn) {
	go func() {
		head := ins.headPool.Get().(*Buffer)
		body := ins.bodyPool.Get().(*Buffer)
		generic := ins.factory(ins.ctx, conn, ins.maxIncomingPacket, head.Bytes, body.Bytes)

		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				log.Println("stack", string(debug.Stack()))
			}
			_ = generic.Close()
			ins.headPool.Put(head)
			ins.bodyPool.Put(body)
		}()

		if appErr := ins.handle(generic); appErr != nil {
			//TODO:
		}
	}()
}

func parseAndWrapOP(ops ...AcceptorOptions) AcceptorOptions {
	var op AcceptorOptions
	if len(ops) > 0 {
		op = ops[0]
	}
	if op.MaxIncomingPacket == 0 {
		op.MaxIncomingPacket = MaxIncomingPacket
	}
	return op
}
