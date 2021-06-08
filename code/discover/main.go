package main

import (
	"ddn/code/discover/dht"
	"encoding/hex"
	"net"

	"github.com/lwch/runtime"
)

func main() {
	var bootstrapAddrs []net.UDPAddr

	for _, addr := range []string{
		"router.bittorrent.com:6881",
		"router.utorrent.com:6881",
		"dht.transmissionbt.com:6881",
	} {
		addr, err := net.ResolveUDPAddr("udp", addr)
		runtime.Assert(err)
		bootstrapAddrs = append(bootstrapAddrs, *addr)
	}
	net, err := dht.New(6882, 1000, 10000, bootstrapAddrs)
	runtime.Assert(err)
	defer net.Close()
	const str = "9f292c93eb0dbdd7ff7a4aa551aaa1ea7cafe004" // debian-10.9.0-amd64-netinst.iso
	hash, err := hex.DecodeString(str)
	runtime.Assert(err)
	var h dht.Hash
	copy(h[:], hash)
	net.Get(h)
	<-make(chan int)
}
