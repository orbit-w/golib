package agent_stream

import (
	"fmt"
	"strings"
	"sync"
)

/*
   @Author: orbit-w
   @File: error
   @2024 4月 周日 23:32
*/

const (
	ErrHeadServe = "agent_stream server serve"
)

type Error struct {
	head string
	text string
}

var builderPool = sync.Pool{New: func() any {
	return &strings.Builder{}
}}

func (e *Error) Error() string {
	w := builderPool.Get().(*strings.Builder)
	defer func() {
		w.Reset()
		builderPool.Put(w)
	}()
	w.WriteString("[")
	w.WriteString(e.head)
	w.WriteString("]: ")
	w.WriteString(e.text)
	return w.String()
}

func New(head string, text string) error {
	return &Error{head: head, text: text}
}

func NewF(head string, format string, args ...interface{}) error {
	return &Error{head: head, text: fmt.Sprintf(format, args...)}
}
