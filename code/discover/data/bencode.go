package data

import "github.com/lwch/bencode"

// ReqType request type
type ReqType string

const (
	// TypePing ping
	TypePing ReqType = "ping"
	// TypeFindNode find_node
	TypeFindNode ReqType = "find_node"
	// TypeGetPeers get_peers
	TypeGetPeers ReqType = "get_peers"
	// TypeAnnouncePeer announce_peer
	TypeAnnouncePeer ReqType = "announce_peer"
)

// Hdr bencode header
type Hdr struct {
	Transaction string `bencode:"t"`
	Type        string `bencode:"y"`
}

func newHdr(t string) Hdr {
	return Hdr{
		Transaction: Rand(16),
		Type:        t,
	}
}

// IsRequest is request packet
func (h Hdr) IsRequest() bool {
	return h.Type == request
}

// IsResponse is response packet
func (h Hdr) IsResponse() bool {
	return h.Type == response
}

// ParseReqType parse request type
func ParseReqType(data []byte) ReqType {
	var t struct {
		Type ReqType `bencode:"q"`
	}
	bencode.Decode(data, &t)
	return t.Type
}
