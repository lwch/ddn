package dht

import (
	"container/list"
	"sync"
)

type reqList struct {
	sync.Mutex
	data list.List
	node *list.Element
}

func newReqList() *reqList {
	return &reqList{}
}

func (l *reqList) push(hash Hash) {
	l.Lock()
	n := l.data.PushBack(hash)
	if l.node == nil {
		l.node = n
	}
	l.Unlock()
}

func (l *reqList) next() Hash {
	l.Lock()
	defer l.Unlock()
	if l.node != nil {
		hash := l.node.Value.(Hash)
		l.node = l.node.Next()
		if l.node == nil {
			l.node = l.data.Front()
		}
		return hash
	}
	return emptyHash
}
