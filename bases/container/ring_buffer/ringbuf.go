package ring_buffer

import (
	"errors"
	"golib/bases/misc/math"
)

// ErrIsEmpty 缓冲区为空
var ErrIsEmpty = errors.New("ring buffer is empty")

// RingBuffer 自动扩容循环缓冲区
type RingBuffer[V any] struct {
	buf      []V
	zero     V
	initSize int
	size     int
	front    int
	rear     int
}

// New Automatically correct the parameter to take the square of the smallest 2 upwards
func New[V any](size int) *RingBuffer[V] {
	//向上取最小的2的平方
	iSize := math.PowerOf2(size)
	return &RingBuffer[V]{
		buf:      make([]V, iSize),
		initSize: iSize,
		size:     iSize,
	}
}

func (ins *RingBuffer[V]) Peek() (V, error) {
	if ins.IsEmpty() {
		return ins.zero, ErrIsEmpty
	}

	v := ins.buf[ins.rear]
	return v, nil
}

func (ins *RingBuffer[V]) Pop() (V, error) {
	if ins.IsEmpty() {
		return ins.zero, ErrIsEmpty
	}
	v := ins.buf[ins.front]
	mask := ins.size - 1
	ins.front = (ins.front + 1) & mask
	return v, nil
}

func (ins *RingBuffer[V]) Push(value V) {
	if ins.IsFull() {
		ins.grow()
	}

	ins.buf[ins.rear] = value
	mask := ins.size - 1
	ins.rear = (ins.rear + 1) & mask
}

func (ins *RingBuffer[V]) Len() int {
	if ins.front == ins.rear {
		return 0
	}
	if ins.rear > ins.front {
		return ins.rear - ins.front
	}

	return ins.size - ins.front + ins.rear
}

func (ins *RingBuffer[V]) Capacity() int {
	return ins.size
}

func (ins *RingBuffer[V]) IsFull() bool {
	mask := ins.size - 1
	return (ins.rear+1)&mask == ins.front
}

func (ins *RingBuffer[V]) IsEmpty() bool {
	return ins.front == ins.rear
}

func (ins *RingBuffer[V]) Reset() {
	ins.front = 0
	ins.rear = 0
	if ins.size > ins.initSize {
		ins.buf = make([]V, ins.initSize)
		ins.size = ins.initSize
	}
}

func (ins *RingBuffer[V]) grow() {
	size := ins.size << 1

	buf := make([]V, size)
	if ins.rear < ins.front {
		copied := copy(buf[0:], ins.buf[ins.front:])
		copy(buf[copied:], ins.buf[0:ins.rear])
	} else {
		copy(buf[0:], ins.buf[ins.front:])
	}
	ins.rear = ins.size - 1
	ins.front = 0
	ins.size = size
	ins.buf = buf
}
