package goswift

import (
	"time"

	"github.com/leoantony72/goswift/expiry"
)

func sweaper(c *cache, h *expiry.Heap) {
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
