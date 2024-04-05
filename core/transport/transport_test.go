package transport

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

var (
	ServeOnce sync.Once
)

func Test_Transport(t *testing.T) {
	host := "127.0.0.1:6800"
	ServeTest(t, host)

	conn := DialWithOps(host, &DialOption{
		RemoteNodeId:  "node_0",
		CurrentNodeId: "node_1",
	})
	defer func() {
		_ = conn.Close()
	}()

	go func() {
		for {
			in, err := conn.Recv()
			if err != nil {
				if IsCancelError(err) || errors.Is(err, io.EOF) {
					log.Println("EOF")
				} else {
					log.Println("Recv failed: ", err.Error())
				}
				break
			}
			log.Println("recv response: ", in[0])
		}
	}()

	w := []byte{1}
	_ = conn.Write(w)

	time.Sleep(time.Second * 10)
}

func ServeTest(t TestingT, host string) {
	ServeOnce.Do(func() {
		listener, err := net.Listen("tcp", host)
		assert.NoError(t, err)
		log.Println("start serve...")

		Serve("tcp", listener, func(conn IServerConn) {
			for {
				in, err := conn.Recv()
				if err != nil {
					if IsClosedConnError(err) {
						break
					}

					if IsCancelError(err) || errors.Is(err, io.EOF) {
						break
					}

					log.Println("conn read stream failed: ", err.Error())
					break
				}
				//log.Println("receive message from client: ", in.Data()[0])
				if err = conn.Send(in); err != nil {
					log.Println("server response failed: ", err.Error())
				}
			}
		})
	})
}

// TestingT is an interface wrapper around *testing.T
type TestingT interface {
	Errorf(format string, args ...interface{})
}
