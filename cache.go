package goswift

import (
	"errors"
	"goswift/expiry"
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
	heapInit := expiry.Init()
	cache := &Cache{Data: datamap, length: 0, heap: heapInit}
	go Sweaper(cache, heapInit)
	return cache
}

func (c *Cache) Exists(key string) bool {
	c.mu.Lock()
	_, ok := c.Data[key]
	c.mu.Unlock()
	return ok
}

// Adds an element to Hash Set, If exp is provided add the
// Node to the Heap with Key and expiration time(int64).
// If exp == 0, Item Never expires, thus it isn't added
// In the Heap
func (c *Cache) Set(key string, exp int, val interface{}) {
	var node *expiry.Node
	if exp == 0 {
		data := &DataHolder{val: val}
		c.mu.Lock()
		c.Data[key] = data
		c.mu.Unlock()
		return
	}
	exp = exp / 1000
	expTime := time.Now().Add(time.Second * time.Duration(exp)).Unix()
	c.mu.Lock()
	node = c.heap.Insert(key, expTime)
	data := &DataHolder{val: val, expiry: node}
	c.Data[key] = data
	c.mu.Unlock()
}

// If exp is not nil, check if the element has expired or not
// Removes the element from the cache if expired, It does not remove
// The Node from the Heap, which will handled by the Sweaper.

// @This must be improved, So that deleted keys does'nt stay in the Heap.
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	data, ok := c.Data[key]
	c.mu.Unlock()
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if data.expiry != nil {
		if data.expiry.Expiry > time.Now().Unix() {
			return data, nil
		}
		c.mu.Lock()
		delete(c.Data, key)
		c.mu.Unlock()
		return nil, errors.New("key does not exist")
	}
	return data.val, nil
}

// Delete's the Item from the Cache

// @This must be improved, So that deleted keys does'nt stay in the Heap.
func (c *Cache) Del(key string) {
	c.mu.Lock()
	delete(c.Data, key)
	c.mu.Unlock()
}

// Set a new value for the key only if it already exist.
// New data will expire at the same time as the prev Key.
func (c *Cache) Update(key string, val interface{}) error {
	if !c.Exists(key) {
		return errors.New("key not present")
	}

	c.mu.Lock()
	e := c.Data[key].expiry
	data := &DataHolder{val: val, expiry: e}
	c.Data[key] = data
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
	val := c.Data[key]
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

func (c *Cache) DeleteExpiredKeys() {
	c.mu.Lock()
	hl := len(c.heap.Data)
	c.mu.Unlock()
	if hl == 0 {
		return
	}
	c.mu.Lock()
	node := c.heap.Data[0]
	c.mu.Unlock()
	// fmt.Println("Node", node)
	if time.Now().Unix() > node.Expiry {
		c.mu.Lock()
		delete(c.Data, node.Key)
		c.heap.Extract()
		c.mu.Unlock()
	}
}
