package dht

import (
	"bytes"
	"container/list"
	"net"
	"sync"

	"github.com/lwch/logging"
)

type bucket struct {
	sync.RWMutex
	prefix Hash
	nodes  *list.List
	leaf   [2]*bucket
	bits   int
}

func (bk *bucket) isLeaf() bool {
	bk.RLock()
	defer bk.RUnlock()
	return bk.leaf[0] == nil && bk.leaf[1] == nil
}

func (bk *bucket) addNode(n *node, k, maxBits int) bool {
	bk.Lock()
	defer bk.Unlock()
	if bk.exists(n.id) {
		// TODO: update
		return false
	}
	if bk.nodes.Len() >= k {
		loopSplit(bk, k, maxBits)
		target := bk.search(n.id)
		if target.exists(n.id) {
			// TODO: update
			return false
		}
		target.nodes.PushBack(n)
		return true
	}
	bk.nodes.PushBack(n)
	return true
}

func loopSplit(bk *bucket, k, maxBits int) {
	bk.split(maxBits)
	if bk.leaf[0] != nil && bk.leaf[0].nodes.Len() >= k {
		loopSplit(bk.leaf[0], k, maxBits)
	}
	if bk.leaf[1] != nil && bk.leaf[1].nodes.Len() >= k {
		loopSplit(bk.leaf[1], k, maxBits)
	}
}

func (bk *bucket) exists(id Hash) bool {
	for n := bk.nodes.Front(); n != nil; n = n.Next() {
		if bytes.Equal(n.Value.(*node).id[:], id[:]) {
			return true
		}
	}
	return false
}

func (bk *bucket) search(id Hash) *bucket {
	if bk.leaf[0] == nil && bk.leaf[1] == nil {
		return bk
	}
	return bk.leaf[id.bit(bk.bits)].search(id)
}

func (bk *bucket) split(maxBits int) {
	if bk.bits >= maxBits {
		return
	}
	var id Hash
	copy(id[:], bk.prefix[:])
	if bk.leaf[0] == nil {
		bk.leaf[0] = newBucket(id, bk.bits+1)
	}
	if bk.leaf[1] == nil {
		bt := bk.bits / 8
		bit := bk.bits % 8
		if bt == 20 {
			var ids []string
			var equals []bool
			for n := bk.nodes.Front(); n != nil; n = n.Next() {
				ids = append(ids, n.Value.(*node).id.String())
				equals = append(equals, bk.equalBits(n.Value.(*node).id))
			}
			logging.Info("overflow: prefix=%s, ids=%v, equals=%v", bk.prefix.String(), ids, equals)
		}
		id[bt] |= 1 << (7 - bit)
		bk.leaf[1] = newBucket(id, bk.bits+1)
	}
	for n := bk.nodes.Front(); n != nil; n = n.Next() {
		node := n.Value.(*node)
		if bk.leaf[0].equalBits(node.id) {
			bk.leaf[0].nodes.PushBack(node)
		} else {
			bk.leaf[1].nodes.PushBack(node)
		}
	}
	bk.nodes = nil
}

func (bk *bucket) equalBits(id Hash) bool {
	bt := bk.bits / 8
	bit := bk.bits % 8
	for i := 0; i < bt; i++ {
		if bk.prefix[i]^id[i] > 0 {
			return false
		}
	}
	a := bk.prefix[bt] >> (8 - bit)
	b := id[bt] >> (8 - bit)
	return a^b <= 0
}

func (bk *bucket) getNodes() []*node {
	bk.RLock()
	defer bk.RUnlock()
	ret := make([]*node, 0, bk.nodes.Len())
	for n := bk.nodes.Front(); n != nil; n = n.Next() {
		ret = append(ret, n.Value.(*node))
	}
	return ret
}

func newBucket(prefix Hash, bits int) *bucket {
	return &bucket{
		prefix: prefix,
		nodes:  list.New(),
		bits:   bits,
	}
}

type table struct {
	sync.RWMutex
	root      *bucket
	addrIndex map[string]*node
	k         int
	size      int
	maxBits   int
}

func bits(n int) int {
	var size int
	for n != 0 {
		size++
		n /= 2
	}
	return size
}

func newTable(k int) *table {
	tb := &table{
		root:      newBucket(emptyHash, 0),
		addrIndex: make(map[string]*node),
		k:         k,
		maxBits:   len(emptyHash)*8 - bits(k),
	}
	return tb
}

func (t *table) add(n *node) bool {
	t.Lock()
	defer t.Unlock()
	next := t.root
	for idx := 0; idx < len(n.id)*8; idx++ {
		if next.isLeaf() {
			ok := next.addNode(n, t.k, t.maxBits)
			if ok {
				t.addrIndex[n.addr.String()] = n
				t.size++
			}
			return ok
		}
		next = next.leaf[n.id.bit(idx)]
	}
	return false
}

func (t *table) remove(n *node) {
	t.Lock()
	defer t.Unlock()
	bk := t.root.search(n.id)
	for nd := bk.nodes.Front(); nd != nil; nd = nd.Next() {
		node := nd.Value.(*node)
		if !node.id.equal(n.id) {
			continue
		}
		delete(t.addrIndex, n.addr.String())
		bk.nodes.Remove(nd)
		t.size--
	}
}

func (t *table) findAddr(addr net.Addr) *node {
	t.RLock()
	data := t.addrIndex[addr.String()]
	t.RUnlock()
	if data == nil {
		return nil
	}
	return data
}

func (t *table) findID(id Hash) *node {
	t.RLock()
	bk := t.root.search(id)
	t.RUnlock()
	for n := bk.nodes.Front(); n != nil; n = n.Next() {
		node := n.Value.(*node)
		if node.id.equal(id) {
			return node
		}
	}
	return nil
}

func (t *table) neighbor(id Hash) []*node {
	t.RLock()
	bk := t.root.search(id)
	t.RUnlock()
	return bk.getNodes()
}
