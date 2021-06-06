package dht

import (
	"net"
)

type node struct {
	id          hashType
	addr        net.UDPAddr
	isBootstrap bool
}

func newNode(id hashType, addr net.UDPAddr) *node {
	return &node{
		id:          id,
		addr:        addr,
		isBootstrap: false,
	}
}

func newBootstrapNode(addr net.UDPAddr) *node {
	return &node{
		id:          randID(),
		addr:        addr,
		isBootstrap: true,
	}
}
