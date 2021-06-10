package dht

import (
	"bytes"
	"context"
	"ddn/code/discover/data"
	"math/rand"
	"net"
	"time"

	"github.com/lwch/bencode"
	"github.com/lwch/logging"
)

type pkt struct {
	data []byte
	addr net.Addr
}

type DHT struct {
	listen *net.UDPConn
	tb     *table
	tx     *txMgr
	local  Hash
	chRead chan pkt

	// runtime
	ctx    context.Context
	cancel context.CancelFunc
	list   *reqList
}

func New(port uint16, addrs []net.UDPAddr) (*DHT, error) {
	dht := &DHT{
		tb:     newTable(8),
		tx:     newTXMgr(30 * time.Second),
		local:  data.RandID(),
		chRead: make(chan pkt, 1000),
		list:   newReqList(),
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
	dht.tx.close()
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
			dht.next()
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
	node := dht.tb.findAddr(addr)
	if node == nil {
		var hdr data.Hdr
		err := bencode.Decode(buf, &hdr)
		if err != nil {
			return
		}
		if hdr.IsRequest() {
			var req struct {
				data.Hdr
				Data struct {
					ID [20]byte `bencode:"id"`
				} `bencode:"a"`
			}
			err = bencode.Decode(buf, &req)
			if err != nil {
				return
			}
			if bytes.Equal(req.Data.ID[:], emptyHash[:]) {
				return
			}
			node = dht.tb.findID(req.Data.ID)
			if node == nil {
				node = newNode(dht, req.Data.ID, *addr.(*net.UDPAddr))
				dht.tb.add(node)
			}
		}
	}
	if node == nil {
		return
	}
	node.onRecv(buf)
}

func (dht *DHT) next() {
	hash := dht.list.next()
	var cnt int
	for cnt < 100 {
		var id Hash
		rand.Read(id[:])
		nodes := dht.tb.neighbor(id)
		for _, node := range nodes {
			node.sendGet(hash)
		}
		cnt += len(nodes)
	}
}

func (dht *DHT) Nodes() int {
	return dht.tb.size
}
