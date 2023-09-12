package expiry

// import "gowi/cache"

// @Max heap
type Heap struct {
	data []*Node
}

type Node struct {
	Kptr   string
	Expiry int64
}

func Init() *Heap {
	return &Heap{}
}

func (h *Heap) Insert(expiry int64, key string) *Node {
	node := &Node{Kptr: key, Expiry: expiry}
	h.data = append(h.data, node)
	h.MaxHeapifyUp(len(h.data) - 1)
	return node
}

func (h *Heap) MaxHeapifyUp(i int) {
	for h.data[parent(i)].Expiry < h.data[i].Expiry {
		h.Swap(parent(i), i)
		i = parent(i)
	}
}

func (h *Heap) Swap(i1, i2 int) {
	h.data[i1], h.data[i2] = h.data[i2], h.data[i1]
}

func parent(index int) int {
	return (index - 1) / 2
}

// func leftchild(index int) int {
// 	return 2*index + 1
// }

// func rightchild(index int) int {
// 	return 2*index + 2
// }
