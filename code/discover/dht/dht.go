package dht

import (
	"context"
	"net"
	"time"
)

type DHT struct {
	listen   *net.UDPConn
	tb       *table
	minNodes int

	// runtime
	ctx    context.Context
	cancel context.CancelFunc
}

func New(port uint16, minNodes, maxNodes int) (*DHT, error) {
	dht := &DHT{
		tb:       newTable(8, maxNodes),
		minNodes: minNodes,
	}
	var err error
	dht.listen, err = net.ListenUDP("udp", &net.UDPAddr{
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	dht.ctx, dht.cancel = context.WithCancel(context.Background())
	go dht.recv()
	return dht, nil
}

func (dht *DHT) Close() {
	dht.listen.Close()
	dht.cancel()
}

func (dht *DHT) recv() {
	buf := make([]byte, 65535)
	for {
		select {
		case <-dht.ctx.Done():
			return
		default:
		}
		dht.listen.SetReadDeadline(time.Now().Add(time.Second))
		n, addr, err := dht.listen.ReadFrom(buf)
		if err != nil {
			continue
		}
	}
}
