package expiry

import "testing"

type tempHeap struct {
	key    string
	expiry int64
}

func TestHeapExpiry(t *testing.T) {
	h := &Heap{}
	buildHeap(h)

	ExpectedValue := []int{1, 2, 6, 12, 14, 23, 123, 436}
	//Check heap values
	t.Run("HeapSortTest", func(t *testing.T) {
		for i := 0; i < len(ExpectedValue); i++ {
			val := Ex(h)
			if int(val.Expiry) != ExpectedValue[i] {
				t.Errorf("Expected Value: %d, Gotten : %d", ExpectedValue[i], val.Expiry)
				return
			}
			// fmt.Println(val.Expiry, ExpectedValue[i])
		}
	})
}

func TestExtract(t *testing.T) {
	h := &Heap{}

	_, err := h.Extract()
	if err != nil {
		if err.Error() != "no elements in the Heap" {
			t.Errorf("expected Error: `no elements in the Heap`, Gotten: %s ", err)
			return
		}
		return
	}
	t.Errorf("expected Error: `no elements in the Heap`, Gotten: %s ", err)

}

func TestHeapInit(t *testing.T) {
	h := Init()

	if h == nil {
		t.Errorf("expected type *Heap gotten nil")
		return
	}
	buildHeap(h)
}

func TestRemove(t *testing.T) {
	h := &Heap{}

	buildHeap(h)

	h.Remove(2, 4)

	h.Remove(1, 1)
}

func Ex(h *Heap) *Node {
	c, _ := h.Extract()
	return c
}

func buildHeap(h *Heap) {
	buildheap := []tempHeap{{"t1", 23}, {"t2", 12}, {"t3", 123}, {"t4", 436}, {"t5", 2}, {"t6", 14}, {"t7", 1}, {"t7", 6}}

	for _, v := range buildheap {
		h.Insert(v.key, v.expiry)
	}
}
