package dht

import (
	"ddn/code/discover/data"
	"fmt"
	"net"

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

func (n *node) errInfo(err error) string {
	return fmt.Sprintf("; id=%s, addr=%s, err=%v",
		n.id.String(), n.addr.String(), err)
}
