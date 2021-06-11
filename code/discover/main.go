package main

import (
	"ddn/code/discover/dht"
	"encoding/hex"
	"net"
	"time"

	"github.com/lwch/logging"
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
	net, err := dht.New(6882, bootstrapAddrs)
	runtime.Assert(err)
	defer net.Close()
	const str = "9f292c93eb0dbdd7ff7a4aa551aaa1ea7cafe004" // debian-10.9.0-amd64-netinst.iso
	hash, err := hex.DecodeString(str)
	runtime.Assert(err)
	var h dht.Hash
	copy(h[:], hash)
	net.Get(h)
	go func() {
		for {
			time.Sleep(time.Minute)
			logging.Info("avg get: %.04f", net.AvgGet())
		}
	}()
	for {
		time.Sleep(time.Second)
		logging.Info("%d nodes", net.Nodes())
	}
}
