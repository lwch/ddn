package data

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
