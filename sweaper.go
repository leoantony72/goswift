package goswift

import (
	"time"

	"github.com/leoantony72/goswift/expiry"
)

// DeleteExpiredKeys func is called every 3 seconds
// which will remove the expired values ffrom the heap and cache
func sweaper(c *Cache, h *expiry.Heap) {
	interval := 3 * time.Second
	// fmt.Println(interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.DeleteExpiredKeys()
		}
	}

}

// An abstraction for testing DeleteExpiredKeys() func
func testDeleteExpiredKeys(c *Cache) {
	c.DeleteExpiredKeys()
}

// Deletes the node if it's expired from both
// Cache(Hash Table) and heap
func (c *Cache) DeleteExpiredKeys() {
	c.mu.Lock()
	l := len(c.Data)
	n := (10 * l) / 100
	if l <= 500 {
		n = 500
	}
	// fmt.Println("N IS ITER: ", n)

	for i := 0; i < n; i++ {
		hl := len(c.heap.Data)
		// fmt.Println(hl)
		if hl == 0 {
			c.mu.Unlock()
			return
		}
		node := c.heap.Data[0]
		// fmt.Println(node)
		if time.Now().Unix() > node.Expiry {
			delete(c.Data, node.Key)
			_, err := c.heap.Extract()
			// fmt.Println(node)
			if err != nil {
				c.mu.Unlock()
				return
			}
			// fmt.Println(c.heap.Data)
		}
	}
	c.mu.Unlock()

}
