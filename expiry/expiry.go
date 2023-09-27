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
	Index  int
	Key    string
	Expiry int64
}

func Init() *Heap {
	return &Heap{}
}

// Takes a key and expiry > 0, Node is added in the last index
// MinheapifyUp func is runned to swap the child node to
// parent node if the exoiry is less than the parent node
func (h *Heap) Insert(key string, expiry int64) *Node {
	node := &Node{Key: key, Expiry: expiry}
	h.Data = append(h.Data, node)
	h.minHeapifyUp(len(h.Data)-1, node)
	return node
}

// swap the child node to parent node if expiry is less than
// of the parent node
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

// Extract function retrives the root node by swapping
// root node with the last node and removing it from the last index
// Last node is then heapifyed Down
func (h *Heap) Extract() (*Node, error) {
	// h.mu.Lock()
	length := len(h.Data)
	if length == 0 {
		// h.mu.Unlock()
		return nil, errors.New("no elements in the Heap")
	}
	node := h.Data[0]
	h.Data[0] = h.Data[length-1]
	h.Data[0].Index = 0
	h.Data = h.Data[:length-1]
	h.minHeapifyDown(0, node)
	// h.mu.Unlock()
	return node, nil
}

// Swaps the parent node with the child node if the child node
// expiry is less than parent node
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

// Remove takes node Index(index of the node to be removed)
// and the last index in the heap, changes the node.Index value
// then swaps them. Finally minHeapifyDown is called to make it
// a valid MinHeap
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

// swap func takes the index of the elements to be swapped
// and then swaps them
func (h *Heap) swap(i1, i2 int) {
	// h.mu.Lock()
	h.Data[i1], h.Data[i2] = h.Data[i2], h.Data[i1]
	// h.mu.Unlock()
}

//formula to get the parent node of a node
func parent(index int) int {
	return (index - 1) / 2
}

//formula to get the left child of a node
func leftchild(index int) int {
	return 2*index + 1
}

//formula to get the right child of a node
func rightchild(index int) int {
	return 2*index + 2
}
