package dht

import (
	"bytes"
	"fmt"
)

type hashType [20]byte

var emptyHash hashType

func (hash hashType) String() string {
	return fmt.Sprintf("%x", [20]byte(hash))
}

func (hash hashType) raw() [20]byte {
	return [20]byte(hash)
}

func (hash hashType) equal(h hashType) bool {
	a := hash.raw()
	b := h.raw()
	return bytes.Equal(a[:], b[:])
}

func (hash hashType) bit(n int) byte {
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
