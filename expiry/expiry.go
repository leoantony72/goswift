package expiry

import (
	"errors"
	"fmt"
	"sync"
)

// @Min Heap
type Heap struct {
	Data []*Node
	mu   sync.RWMutex
}

type Node struct {
	Index  int
	Key    string
	Expiry int64
}

func Init() *Heap {
	return &Heap{}
}

func (h *Heap) Insert(key string, expiry int64) *Node {
	// h.mu.Lock()
	// if len(h.Data) ==
	// l := len(h.Data)
	// if l == 0 {
	// 	l += 0
	// } else {
	// 	l -= 1
	// }
	fmt.Println("hlength", len(h.Data))
	node := &Node{Key: key, Expiry: expiry}
	h.Data = append(h.Data, node)
	h.minHeapifyUp(len(h.Data)-1, node)
	// h.mu.Unlock()
	return node
}

func (h *Heap) minHeapifyUp(i int, node *Node) {
	for i > 0 && h.Data[parent(i)].Expiry > h.Data[i].Expiry {
		h.Data[parent(i)].Index = i
		// h.mu.Unlock()
		node.Index = parent(i)
		h.swap(parent(i), i)
		// h.mu.Lock()
		i = parent(i)
	}
	node.Index = i
}

func (h *Heap) Extract() (*Node, error) {
	// h.mu.Lock()
	length := len(h.Data)
	if length == 0 {
		h.mu.Unlock()
		return nil, errors.New("no elements in the Heap")
	}
	node := h.Data[0]
	h.Data[0] = h.Data[length-1]
	h.Data = h.Data[:length-1]
	h.minHeapifyDown(0, node)
	// h.mu.Unlock()
	return node, nil
}

func (h *Heap) minHeapifyDown(i int, node *Node) {
	// h.mu.RLock()
	lastIndex := len(h.Data) - 1
	// h.mu.RUnlock()

	childToCompare := 0
	l, r := leftchild(i), rightchild(i)
	for l <= lastIndex {
		// h.mu.Lock()
		if l == lastIndex {
			childToCompare = l
		} else if h.Data[l].Expiry < h.Data[r].Expiry {
			childToCompare = l
		} else {
			childToCompare = r
		}

		if h.Data[childToCompare].Expiry < h.Data[i].Expiry {
			h.Data[childToCompare].Index = i
			h.Data[i].Index = childToCompare
			h.swap(childToCompare, i)
			i = childToCompare
			l, r = leftchild(i), rightchild(i)
		} else {
			// h.mu.Unlock()
			return
		}

		// h.mu.Unlock()
	}

}

func (h *Heap) Remove(nindex, lindex int) {
	if nindex == lindex {
		h.Data = h.Data[:len(h.Data)-1]
		return
	}
	// h.mu.Lock()
	h.Data[nindex].Index = lindex
	h.Data[lindex].Index = nindex
	h.swap(nindex, lindex)
	h.Data = h.Data[:len(h.Data)-1]
	h.minHeapifyDown(nindex, h.Data[nindex])
	// h.mu.Unlock()
}

func (h *Heap) swap(i1, i2 int) {
	// h.mu.Lock()
	h.Data[i1], h.Data[i2] = h.Data[i2], h.Data[i1]
	// h.mu.Unlock()
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
