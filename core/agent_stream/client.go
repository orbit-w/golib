package agent_stream

import (
	"context"
	"github.com/orbit-w/golib/core/network"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"sync"
	"sync/atomic"
	"time"
)

type SClient interface {
	Stream() (IStream, error)
}

type Client struct {
	addr     string
	state    atomic.Uint32
	xClient  client.XClient
	conf     *ClientConfig
	headPool *sync.Pool
	bodyPool *sync.Pool
}

func NewClient(addr string) *Client {
	conf := DefaultClientConfig()
	return NewClientByConfig(addr, conf)
}

func NewClientByConfig(addr string, conf *ClientConfig) *Client {
	c := &Client{
		addr:     addr,
		conf:     conf,
		headPool: network.NewBufferPool(network.HeadLen),
		bodyPool: network.NewBufferPool(conf.MaxIncomingPacket),
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

	conn, err := c.xClient.Stream(ctx, make(map[string]string))
	if err != nil {
		return nil, err
	}
	stream := NewAgentStream(conn, conf.MaxIncomingPacket, conf.IsGzip, conf.WriteTimeout, conf.ReadTimeout)
	go func() {
		var (
			head = c.headPool.Get().(*network.Buffer)
			body = c.bodyPool.Get().(*network.Buffer)
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
		conf.MaxIncomingPacket = network.MaxIncomingPacket
	}
	if conf.WriteTimeout == 0 {
		conf.WriteTimeout = WriteTimeout
	}
	if conf.ReadTimeout == 0 {
		conf.ReadTimeout = ReadTimeout
	}
	if conf.MaxIncomingPacket == 0 {
		conf.MaxIncomingPacket = network.MaxIncomingPacket
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
		MaxIncomingPacket: network.MaxIncomingPacket,
		IsGzip:            false,
		ReadTimeout:       ReadTimeout,
		DialTimeout:       DialTimeout,
	}
}
