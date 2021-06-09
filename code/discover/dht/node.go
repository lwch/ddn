package dht

import (
	"ddn/code/discover/data"
	"fmt"
	"net"

	"github.com/lwch/bencode"
	"github.com/lwch/logging"
)

type node struct {
	dht         *DHT
	id          Hash
	addr        net.UDPAddr
	isBootstrap bool
}

func newNode(dht *DHT, id Hash, addr net.UDPAddr) *node {
	return &node{
		dht:         dht,
		id:          id,
		addr:        addr,
		isBootstrap: false,
	}
}

func newBootstrapNode(dht *DHT, addr net.UDPAddr) *node {
	return &node{
		dht:         dht,
		id:          randID(),
		addr:        addr,
		isBootstrap: true,
	}
}

func (n *node) sendDiscovery(id Hash) {
	pkt, tx, err := data.FindReq(n.dht.local, id)
	if err != nil {
		logging.Error("build find_node packet failed" + n.errInfo(err))
		return
	}
	_, err = n.dht.listen.WriteTo(pkt, &n.addr)
	if err != nil {
		logging.Error("send find_node packet failed" + n.errInfo(err))
		return
	}
	n.dht.tx.add(tx, data.TypeFindNode, emptyHash, id)
}

func (n *node) sendGet(hash Hash) {
	buf, tx, err := data.GetPeers(n.dht.local, hash)
	if err != nil {
		logging.Error("build get_peers packet failed" + n.errInfo(err))
		return
	}
	_, err = n.dht.listen.WriteTo(buf, &n.addr)
	if err != nil {
		logging.Error("send get_peers packet failed" + n.errInfo(err))
		return
	}
	n.dht.tx.add(tx, data.TypeGetPeers, hash, emptyHash)
}

func (n *node) onRecv(buf []byte) {
	var hdr data.Hdr
	err := bencode.Decode(buf, &hdr)
	if err != nil {
		// TODO: log
		return
	}
	switch {
	case hdr.IsRequest():
		n.handleRequest(buf)
	case hdr.IsResponse():
		n.handleResponse(buf, hdr.Transaction)
	}
}

func (n *node) handleRequest(buf []byte) {
	var req struct {
		Data struct {
			ID [20]byte `bencode:"id"`
		} `bencode:"a"`
	}
	err := bencode.Decode(buf, &req)
	if err != nil {
		logging.Error("decode request failed" + n.errInfo(err))
		return
	}
	if !n.id.equal(req.Data.ID) {
		n.dht.tb.remove(n)
		return
	}
	switch data.ParseReqType(buf) {
	case data.TypePing:
		n.onPing(buf)
	case data.TypeFindNode:
		n.onFindNode(buf)
	case data.TypeGetPeers:
		n.onGetPeers(buf)
	case data.TypeAnnouncePeer:
		n.onAnnouncePeer(buf)
	}
}

func (n *node) handleResponse(buf []byte, tx string) {
	txr := n.dht.tx.find(tx)
	if txr == nil {
		return
	}
	switch txr.t {
	case data.TypePing:
	case data.TypeFindNode:
		n.onFindNodeResp(buf)
	case data.TypeGetPeers:
		n.onGetPeersResp(buf, txr.hash)
	}
}

func (n *node) errInfo(err error) string {
	return fmt.Sprintf("; id=%s, addr=%s, err=%v",
		n.id.String(), n.addr.String(), err)
}
