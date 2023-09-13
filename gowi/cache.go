package gowi

import (
	"errors"
	"fmt"
	"gowiz/gowi/expiry"
	"sync"
	"time"
)

type Cache struct {
	Data   map[string]*DataHolder
	length int
	heap   *expiry.Heap
	mu     sync.Mutex
}

type DataHolder struct {
	val    interface{}
	expiry *expiry.Node
}

func NewCache() *Cache {
	datamap := make(map[string]*DataHolder)
	heap := expiry.Init()
	cache := &Cache{Data: datamap, length: 0, heap: heap}
	go Sweaper(cache, heap)
	return cache
}

func (c *Cache) Exists(key string) bool {
	c.mu.Lock()
	_, ok := c.Data[key]
	c.mu.Unlock()
	if ok {
		return true
	}
	return false
}

func (c *Cache) Set(key string, exp int, val interface{}) {
	var node *expiry.Node
	if exp != 0 {
		exp = exp / 1000
		expTime := time.Now().Add(time.Second * time.Duration(exp)).Unix()
		c.mu.Lock()
		node = c.heap.Insert(expTime, key)
		datamap := &DataHolder{val: val, expiry: node}
		c.Data[key] = datamap
		c.mu.Unlock()
		return
	}
	datamap := &DataHolder{val: val}
	c.mu.Lock()
	c.Data[key] = datamap
	c.mu.Unlock()
	// fmt.Println(c.Data[key])
}

func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	val, ok := c.Data[key]
	c.mu.Unlock()
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if val.expiry != nil {
		// fmt.Println(val.expiry.Expiry, time.Now().Unix())
		if val.expiry.Expiry > time.Now().Unix() {
			return val, nil
		}
		c.mu.Lock()
		delete(c.Data, key)
		c.mu.Unlock()
		return nil, errors.New("key does not exist")
	}
	return val.val, nil
}

func (c *Cache) Del(key string) {
	c.mu.Lock()
	delete(c.Data, key)
	c.mu.Unlock()
}

func (c *Cache) Update(key string, val interface{}) error {
	//check if key is present
	// if _, ok := c.data[key]; !ok {
	// 	return errors.New("key not present")
	// }

	if !c.Exists(key) {
		return errors.New("key not present")
	}
	datamap := &DataHolder{val: val}
	c.mu.Lock()
	c.Data[key] = datamap
	c.mu.Unlock()
	return nil
}

func (c *Cache) Hset(key, field string, value interface{}) {

	c.mu.Lock()
	if _, exists := c.Data[key]; !exists {
		c.Data[key] = &DataHolder{}
		c.Data[key].val = make(map[string]interface{})
	}

	hash := c.Data[key].val.(map[string]interface{})
	hash[field] = value
	c.mu.Unlock()
}

func (c *Cache) HGet(key, field string) (interface{}, error) {
	if !c.Exists(key) {
		return nil, errors.New("key not present")
	}
	c.mu.Lock()
	val, _ := c.Data[key]
	c.mu.Unlock()
	if mpval, ok := val.val.(map[string]interface{}); ok {
		if data, ok := mpval[field]; ok {
			return data, nil
		}
		return nil, errors.New("key not present")
	}
	return nil, errors.New("not a Hash value/table")

}

func (c *Cache) HGetAll(key string) (map[string]interface{}, error) {
	c.mu.Lock()
	data, ok := c.Data[key]
	c.mu.Unlock()
	if ok {

		if mpdata, oks := data.val.(map[string]interface{}); oks {
			return mpdata, nil
		}
		return nil, errors.New("not a Hash value/table")
	}
	return nil, errors.New("key not present")
}

func Sweaper(c *Cache, h *expiry.Heap) {
	interval := 2 * time.Second
	// fmt.Println(interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			{
				fmt.Println("concurrency active")
				// fmt.Println(len(h.Data))
				c.mu.Lock()
				heaplength := len(h.Data)
				c.mu.Unlock()
				if heaplength == 0 {
					continue
				}
				lastIndex := heaplength - 1
				c.mu.Lock()
				heapData := h.Data[lastIndex]
				c.mu.Unlock()
				fmt.Println(heapData)
				if heapData.Expiry < time.Now().Unix() {
					c.mu.Lock()
					delete(c.Data, heapData.Kptr)
					h.Data = h.Data[:lastIndex]
					c.mu.Unlock()
				}
				continue
			}
		}
	}

}
