package network

import (
	"context"
	"github.com/orbit-w/golib/bases/misc/utils"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
   @Author: orbit-w
   @File: server
   @2023 11月 周五 17:04
*/

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
	handle            ConnHandle
	bodyPool          *sync.Pool
	headPool          *sync.Pool
}

type AcceptorOptions struct {
	MaxIncomingPacket uint32
	IsGzip            bool
}

func (ins *Server) Serve(p Protocol, listener net.Listener, _handle ConnHandle, ops ...AcceptorOptions) {
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

	ins.headPool = NewBufferPool(HeadLen)
	ins.bodyPool = NewBufferPool(ins.maxIncomingPacket)

	go ins.acceptLoop()
}

// Stop stops the server
// 具有可重入性且线程安全, 这意味着这个方法可以被并发多次调用，而不会影响程序的状态或者产生不可预期的结果
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
	utils.GoRecoverPanic(func() {
		head := ins.headPool.Get().(*Buffer)
		body := ins.bodyPool.Get().(*Buffer)
		defer func() {
			ins.headPool.Put(head)
			ins.bodyPool.Put(body)
		}()

		ins.handle(ins.ctx, conn, ins.maxIncomingPacket, head.Bytes, body.Bytes)
	})
}

func DefaultAcceptorOptions() AcceptorOptions {
	return AcceptorOptions{
		MaxIncomingPacket: MaxIncomingPacket,
		IsGzip:            false,
	}
}

func parseAndWrapOP(ops ...AcceptorOptions) AcceptorOptions {
	var op AcceptorOptions
	if len(ops) > 0 {
		op = ops[0]
		if op.MaxIncomingPacket == 0 {
			op.MaxIncomingPacket = MaxIncomingPacket
		}
	} else {
		op = DefaultAcceptorOptions()
	}

	return op
}
