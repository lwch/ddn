package dht

import (
	"context"
	"ddn/code/discover/data"
	"encoding/hex"
	"net"
	"time"

	"github.com/lwch/logging"
)

type pkt struct {
	data []byte
	addr net.Addr
}

type DHT struct {
	listen   *net.UDPConn
	tb       *table
	tx       *txMgr
	minNodes int
	local    Hash
	chRead   chan pkt

	// runtime
	ctx    context.Context
	cancel context.CancelFunc
	list   *reqList
}

func New(port uint16, minNodes, maxNodes int, addrs []net.UDPAddr) (*DHT, error) {
	dht := &DHT{
		tb:       newTable(8, maxNodes),
		tx:       newTXMgr(30 * time.Second),
		minNodes: minNodes,
		local:    data.RandID(),
		chRead:   make(chan pkt, 1000),
		list:     newReqList(),
	}
	for _, addr := range addrs {
		n := newBootstrapNode(dht, addr)
		dht.tb.add(n)
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
	go dht.handler()
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
		data := make([]byte, n)
		copy(data, buf[:n])
		select {
		case dht.chRead <- pkt{
			data: data,
			addr: addr,
		}:
		default:
			logging.Info("drop packet")
		}
	}
}

func (dht *DHT) handler() {
	tk := time.NewTicker(time.Second)
	for {
		select {
		case pkt := <-dht.chRead:
			dht.handleData(pkt.addr, pkt.data)
		case <-tk.C:
			if dht.tb.size < dht.minNodes {
				dht.nextGet()
			} else if dht.tx.size() == 0 {
				dht.nextGet()
			}
		case <-dht.ctx.Done():
			return
		}
	}
}

func (dht *DHT) Get(hash Hash) {
	nodes := dht.tb.neighbor(hash)
	for _, node := range nodes {
		node.sendGet(hash)
	}
	dht.list.push(hash)
}

func (dht *DHT) handleData(addr net.Addr, buf []byte) {
	logging.Info("addr=%s\n%s", addr.String(), hex.Dump(buf))
}

func (dht *DHT) nextGet() {
	hash := dht.list.next()
	nodes := dht.tb.neighbor(hash)
	for _, node := range nodes {
		node.sendGet(hash)
	}
}
