package dht

import (
	"bytes"
	"fmt"
)

type Hash [20]byte

var emptyHash Hash

func (hash Hash) String() string {
	return fmt.Sprintf("%x", [20]byte(hash))
}

func (hash Hash) raw() [20]byte {
	return [20]byte(hash)
}

func (hash Hash) equal(h Hash) bool {
	a := hash.raw()
	b := h.raw()
	return bytes.Equal(a[:], b[:])
}

func (hash Hash) bit(n int) byte {
	if n < 0 || n >= len(hash)*8 {
		panic(fmt.Errorf("out of range 0~%d[%d]", len(hash)*8-1, n))
	}
	bt := n / 8
	bit := n % 8
	if bit > 0 {
		return (hash[bt] >> (7 - bit)) & 1
	}
	return hash[bt] >> 7
}
