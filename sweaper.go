package gowiz

import (
	"gowiz/expiry"
	"time"
)

func Sweaper(c *Cache, h *expiry.Heap) {
	interval := 2 * time.Second
	// fmt.Println(interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go c.DeleteExpiredKeys()
		}
	}

}
