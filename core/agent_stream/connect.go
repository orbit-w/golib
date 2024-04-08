package agent_stream

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

/*
   @Author: orbit-w
   @File: connect
   @2024 4月 周一 22:31
*/

func Dial(addr string) (IStream, error) {
	conf := DefaultConfig()
	return DialByConfig(addr, conf)
}

func DialByConfig(addr string, conf *Config) (IStream, error) {
	parseConfig(conf)

	d, err := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	if err != nil {
		return nil, err
	}
	xClient := client.NewXClient(share.StreamServiceName, client.Failtry, client.RandomSelect, d, client.DefaultOption)
	ctx, cancel := context.WithTimeout(context.Background(), conf.DialTimeout)
	defer cancel()

	conn, err := xClient.Stream(ctx, make(map[string]string))
	if err != nil {
		return nil, err
	}
	stream := NewAgentStream(conn, conf)
	var (
		head = make([]byte, 4)
		body = make([]byte, conf.MaxIncomingPacket)
	)
	go stream.handleLoop(conn, head, body)
	return stream, nil
}
