package transport

import (
	"errors"
	"strings"
)

/*
   @Author: orbit-w
   @File: error
   @2023 11月 周日 14:29
*/

var (
	ErrRpcDisconnected  = errors.New("error_rpc_disconnected")
	ErrRpcDisconnectedP = "error_rpc_disconnected"
	ErrStreamDone       = errors.New("error_the_stream_is_done")
)

func IsErrRpcDisconnected(err error) bool {
	return err != nil && strings.Contains(err.Error(), ErrRpcDisconnectedP)
}

func IsClosedConnError(err error) bool {
	/*
		`use of closed file or network connection` (Go ver > 1.8, internal/pool.ErrClosing)
		`mux: listener closed` (cmux.ErrListenerClosed)
	*/
	return err != nil && strings.Contains(err.Error(), "closed")
}
