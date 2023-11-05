package transport

import (
	"github.com/orbit-w/golib/v1/bases/misc/number_utils"
	"github.com/orbit-w/golib/v1/bases/packet"
	"sync"
)

/*
   @Author: orbit-w
   @File: control_buf
   @2023 11月 周日 17:21
*/

type ControlBuffer struct {
	consumerWaiting bool
	state           int8
	max             uint32
	length          int
	buffer          packet.IPacket
	mu              sync.Mutex
	sw              *SenderWrapper

	ch    chan struct{}
	close chan struct{}
}

func NewControlBuffer(max uint32, _sw *SenderWrapper) *ControlBuffer {
	ins := &ControlBuffer{
		state:           TypeWorking,
		consumerWaiting: false,
		max:             max,
		buffer:          packet.New(),
		mu:              sync.Mutex{},
		ch:              make(chan struct{}, 1),
		close:           make(chan struct{}, 1),
		sw:              _sw,
	}
	go ins.flush()
	ins.Kick()
	return ins
}

func (ins *ControlBuffer) Run(_sw *SenderWrapper) {
	ins.mu.Lock()
	if ins.state == TypeStopped {
		ins.mu.Unlock()
		return
	}
	ins.sw = _sw
	ins.close = make(chan struct{}, 1)
	ins.state = TypeWorking
	ins.mu.Unlock()

	go ins.flush()
	ins.Kick()
}

func (ins *ControlBuffer) Kick() {
	var kick bool
	if ins.state == TypeWorking && ins.consumerWaiting {
		kick = true
		ins.consumerWaiting = false
	}
	if kick {
		select {
		case ins.ch <- struct{}{}:
		default:

		}
	}
}

func (ins *ControlBuffer) Set(buf packet.IPacket) error {
	ins.mu.Lock()
	if ins.state == TypeStopped {
		ins.mu.Unlock()
		return ErrRpcDisconnected
	}
	var kick bool
	ins.length++
	ins.buffer.WriteBytes32(buf.Data())
	if ins.consumerWaiting {
		kick = true
		ins.consumerWaiting = false
	}
	ins.mu.Unlock()
	if kick {
		select {
		case ins.ch <- struct{}{}:
		default:
		}
	}
	return nil
}

//TODO: 立即拒绝流输入是否合理？
func (ins *ControlBuffer) OnClose() {
	ins.mu.Lock()
	defer ins.mu.Unlock()
	ins.state = TypeStopped
	close(ins.close)
}

func (ins *ControlBuffer) flush() {
	defer func() {
		if x := recover(); x != nil {
		}
		ins.safeReturn()
	}()

FLUSH:
	ins.mu.Lock()
	for ins.state == TypeWorking && !ins.isEmpty() {
		w := packet.Writer()
		size := number_utils.Min[int](BatchLimit, ins.length)
		for i := 0; i < size; i++ {
			length, _ := ins.buffer.NextBytesSize32()
			if uint32(w.Len())+uint32(length)+4 > ins.max {
				break
			}
			ins.length++
			data, _ := ins.buffer.ReadBytes32()
			w.WriteBytes32(data)
		}

		for {
			if err := ins.sw.Send(w); err == nil {
				break
			}
		}
	}

	ins.consumerWaiting = true
	ins.buffer.Reset()
	ins.mu.Unlock()
	select {
	case <-ins.ch:
		goto FLUSH
	case <-ins.close:
		return
	}
}

func (ins *ControlBuffer) safeReturn() {
	ins.mu.Lock()
	defer ins.mu.Unlock()
	ins.state = TypeStopped
	ins.sw.OnClose()
	ins.buffer.Reset()
	close(ins.ch)
	ins.buffer.Return()
	ins.buffer = nil
}

func (ins *ControlBuffer) isEmpty() bool {
	return ins.length == 0
}
