package data

import "github.com/lwch/bencode"

const (
	request  = "q"
	response = "r"
	err      = "e"
)

// GetPeersRequest get_peers request
type GetPeersRequest struct {
	Hdr
	Action string `bencode:"q"`
	Data   struct {
		ID   [20]byte `bencode:"id"`
		Hash [20]byte `bencode:"info_hash"`
	} `bencode:"a"`
}

// GetPeers build get_peers request packet
func GetPeers(id, hash [20]byte) ([]byte, string, error) {
	var req GetPeersRequest
	req.Hdr = newHdr(request)
	req.Action = "get_peers"
	req.Data.ID = id
	req.Data.Hash = hash
	data, err := bencode.Encode(req)
	if err != nil {
		return nil, "", err
	}
	return data, req.Hdr.Transaction, nil
}
