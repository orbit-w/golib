package ring_buffer

import (
	"errors"
	"fmt"
	math2 "github.com/orbit-w/golib/bases/misc/math"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer_PowerOf2(t *testing.T) {
	fmt.Println(math2.PowerOf2(math.MaxInt))
	fmt.Println(math.MaxUint32)
	fmt.Println(math2.PowerOf2(-1))

	fmt.Println(1024 << 1)
}

func TestRingBuf(t *testing.T) {
	rb := New[int](10)
	v, err := rb.Pop()
	assert.Error(t, err, ErrIsEmpty)

	var write, read int

	rb.Push(0)
	v, err = rb.Pop()
	assert.NoError(t, err)
	assert.Equal(t, 0, v)
	assert.Equal(t, 1, rb.rear)
	assert.Equal(t, 1, rb.front)
	assert.True(t, rb.IsEmpty())

	for i := 1; i < 10; i++ {
		rb.Push(i)
		write += i
	}
	assert.Equal(t, math2.PowerOf2(10), rb.Capacity())
	assert.Equal(t, 9, rb.Len())

	rb.Push(10)
	write += 10
	assert.Equal(t, math2.PowerOf2(10), rb.Capacity())
	assert.Equal(t, 10, rb.Len())

	for i := 1; i <= 90; i++ {
		rb.Push(i)
		write += i
	}

	assert.Equal(t, 128, rb.Capacity())
	assert.Equal(t, 100, rb.Len())

	for {
		v, err := rb.Pop()
		if errors.Is(err, ErrIsEmpty) {
			break
		}
		read += v
	}

	assert.Equal(t, write, read)
	rb.Reset()
	assert.Equal(t, 16, rb.Capacity())
	assert.Equal(t, 0, rb.Len())
	assert.True(t, rb.IsEmpty())
}
