package expiry

import "testing"

func TestHeapExpiry(t *testing.T) {
	h := &Heap{}
	type th struct {
		key    string
		expiry int64
	}
	buildheap := []th{{"t1", 23}, {"t2", 12}, {"t3", 123}, {"t4", 436}, {"t5", 2}, {"t6", 14}, {"t7", 1}, {"t7", 6}}

	for _, v := range buildheap {
		h.Insert(v.key, v.expiry)
	}

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

func Ex(h *Heap) *Node {
	c, _ := h.Extract()
	return c
}
