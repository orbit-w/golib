package agent_stream

import (
	"context"
	network2 "github.com/orbit-w/golib/modules/net/network"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"sync"
	"sync/atomic"
	"time"
)

type ASClient interface {
	Stream() (IStream, error)
	Close() error
}

type Client struct {
	addr     string
	state    atomic.Uint32
	xClient  client.XClient
	conf     *ClientConfig
	headPool *sync.Pool
	bodyPool *sync.Pool
}

func NewClient(addr string) ASClient {
	conf := DefaultClientConfig()
	return NewClientByConfig(addr, conf)
}

func NewClientByConfig(addr string, conf *ClientConfig) *Client {
	c := &Client{
		addr:     addr,
		conf:     conf,
		headPool: network2.NewBufferPool(network2.HeadLen),
		bodyPool: network2.NewBufferPool(conf.MaxIncomingPacket),
	}
	c.parseConfig(conf)

	d, _ := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	xClient := client.NewXClient(share.StreamServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	c.xClient = xClient
	return c
}

func (c *Client) Stream() (IStream, error) {
	return c.dialByConfig(c.conf)
}

func (c *Client) Close() error {
	if c.state.CompareAndSwap(StateNormal, StateStopped) {
		if c.xClient != nil {
			_ = c.xClient.Close()
		}
	}
	return nil
}

func (c *Client) dialByConfig(conf *ClientConfig) (IStream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), conf.DialTimeout)
	defer cancel()

	conn, err := c.xClient.Stream(ctx, map[string]string{"foo": "bar"})
	if err != nil {
		return nil, err
	}
	stream := NewAgentStream(conn, conf.MaxIncomingPacket, conf.IsGzip, conf.WriteTimeout, conf.ReadTimeout)
	go func() {
		var (
			head = c.headPool.Get().(*network2.Buffer)
			body = c.bodyPool.Get().(*network2.Buffer)
		)
		defer func() {
			c.headPool.Put(head)
			c.bodyPool.Put(body)
		}()
		stream.handleLoop(conn, head.Bytes, body.Bytes)
	}()
	return stream, nil
}

func (c *Client) parseConfig(conf *ClientConfig) {
	if conf.MaxIncomingPacket <= 0 {
		conf.MaxIncomingPacket = network2.MaxIncomingPacket
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = WriteTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = ReadTimeout
	}
	if conf.MaxIncomingPacket == 0 {
		conf.MaxIncomingPacket = network2.MaxIncomingPacket
	}
	if conf.DialTimeout == 0 {
		conf.DialTimeout = DialTimeout
	}
}

type ClientConfig struct {
	MaxIncomingPacket uint32
	IsGzip            bool
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	DialTimeout       time.Duration
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		MaxIncomingPacket: network2.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		DialTimeout:       DialTimeout,
	}
}
