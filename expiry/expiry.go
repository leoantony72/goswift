package expiry

import (
	"errors"
	"sync"
)

// @Min Heap
type Heap struct {
	Data []*Node
	mu   sync.RWMutex
}

type Node struct {
	Key    string
	Expiry int64
}

func Init() *Heap {
	return &Heap{}
}

func (h *Heap) Insert(key string, expiry int64) *Node {
	node := &Node{Key: key, Expiry: expiry}
	h.mu.Lock()
	h.Data = append(h.Data, node)
	h.mu.Unlock()
	h.minHeapifyUp(len(h.Data) - 1)
	return node
}

func (h *Heap) minHeapifyUp(i int) {
	h.mu.RLock()
	for i > 0 && h.Data[parent(i)].Expiry > h.Data[i].Expiry {
		h.swap(parent(i), i)
		i = parent(i)
	}
	h.mu.RUnlock()
}

func (h *Heap) Extract() (*Node, error) {
	h.mu.RLock()
	length := len(h.Data)
	h.mu.RUnlock()
	if length == 0 {
		return nil, errors.New("no elements in the Heap")
	}
	h.mu.Lock()
	node := h.Data[0]
	h.Data[0] = h.Data[length-1]
	h.Data = h.Data[:length-1]
	h.mu.Unlock()
	h.minHeapifyDown(0)
	return node, nil
}

func (h *Heap) minHeapifyDown(i int) {
	h.mu.RLock()
	lastIndex := len(h.Data) - 1
	h.mu.RUnlock()

	childToCompare := 0
	l, r := leftchild(i), rightchild(i)
	for l <= lastIndex {
		h.mu.Lock()
		if l == lastIndex {
			childToCompare = l
		} else if h.Data[l].Expiry < h.Data[r].Expiry {
			childToCompare = l
		} else {
			childToCompare = r
		}

		if h.Data[childToCompare].Expiry < h.Data[i].Expiry {
			h.swap(childToCompare, i)
			i = childToCompare
			l, r = leftchild(i), rightchild(i)
		} else {
			h.mu.Unlock()
			return
		}

		h.mu.Unlock()
	}

}

func (h *Heap) swap(i1, i2 int) {
	h.mu.Lock()
	h.Data[i1], h.Data[i2] = h.Data[i2], h.Data[i1]
	h.mu.Unlock()
}

func parent(index int) int {
	return (index - 1) / 2
}

func leftchild(index int) int {
	return 2*index + 1
}

func rightchild(index int) int {
	return 2*index + 2
}
