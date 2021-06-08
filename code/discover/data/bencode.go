package data

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
