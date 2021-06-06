package main

import (
	"ddn/code/discover/dht"

	"github.com/lwch/runtime"
)

func main() {
	net, err := dht.New(6881, 1000, 10000)
	runtime.Assert(err)
	defer net.Close()
}
