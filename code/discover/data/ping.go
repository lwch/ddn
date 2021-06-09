package data

import "github.com/lwch/bencode"

// PingResponse ping response
type PingResponse struct {
	Hdr
	Response struct {
		ID [20]byte `bencode:"id"`
	} `bencode:"r"`
}

// PingRep build ping response packet
func PingRep(tx string, id [20]byte) ([]byte, error) {
	var rep PingResponse
	rep.Transaction = tx
	rep.Type = response
	rep.Response.ID = id
	return bencode.Encode(rep)
}
