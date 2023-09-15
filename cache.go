package goswift

import (
	"errors"
	"sync"
	"time"

	"github.com/leoantony72/goswift/expiry"
)

type cache struct {
	Data   map[string]*dataHolder
	length int
	heap   *expiry.Heap
	mu     sync.RWMutex
}

type dataHolder struct {
	Value  interface{}
	Expiry *expiry.Node
}

func NewCache() cacheFunction {
	dataMap := make(map[string]*dataHolder)
	heapInit := expiry.Init()
	cache := &cache{Data: dataMap, length: 0, heap: heapInit}
	go sweaper(cache, heapInit)
	return cache
}

func (c *cache) Exists(key string) bool {
	c.mu.RLock()
	_, ok := c.Data[key]
	c.mu.RUnlock()
	return ok
}

// Adds an element to Hash Set, If exp is provided add the
// Node to the Heap with Key and expiration time(int64).
// If exp == 0, Item Never expires, thus it isn't added
// In the Heap
func (c *cache) Set(key string, exp int, val interface{}) {
	var node *expiry.Node
	if exp == 0 {
		data := &dataHolder{Value: val}
		c.mu.Lock()
		c.Data[key] = data
		c.mu.Unlock()
		return
	}
	exp = exp / 1000
	expTime := time.Now().Add(time.Second * time.Duration(exp)).Unix()
	node = c.heap.Insert(key, expTime)
	data := &dataHolder{Value: val, Expiry: node}
	c.mu.Lock()
	c.Data[key] = data
	c.mu.Unlock()
}

// If exp is not nil, check if the element has expired or not
// Removes the element from the cache if expired, It does not remove
// The Node from the Heap, which will handled by the Sweaper.
// @This must be improved, So that deleted keys does'nt stay in the Heap.
func (c *cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	data, ok := c.Data[key]
	c.mu.RUnlock()
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if data.Expiry != nil {
		if data.Expiry.Expiry > time.Now().Unix() {
			return data.Value, nil
		}
		c.mu.Lock()
		delete(c.Data, key)
		c.mu.Unlock()
		return nil, errors.New("key does not exist")
	}
	return data.Value, nil
}

// Delete's the Item from the cache
// @This must be improved, So that deleted keys does'nt stay in the Heap.
func (c *cache) Del(key string) {
	c.mu.Lock()
	delete(c.Data, key)
	c.mu.Unlock()
}

// Set a new value for the key only if it already exist.
// New data will expire at the same time as the prev Key.
func (c *cache) Update(key string, val interface{}) error {
	if !c.Exists(key) {
		return errors.New("key not present")
	}

	c.mu.Lock()
	e := c.Data[key].Expiry
	data := &dataHolder{Value: val, Expiry: e}
	c.Data[key] = data
	c.mu.Unlock()
	return nil
}

func (c *cache) Hset(key, field string, value interface{}) {

	c.mu.Lock()
	if _, exists := c.Data[key]; !exists {
		c.Data[key] = &dataHolder{}
		c.Data[key].Value = make(map[string]interface{})
	}

	hash := c.Data[key].Value.(map[string]interface{})
	hash[field] = value
	c.mu.Unlock()
}

func (c *cache) HGet(key, field string) (interface{}, error) {
	if !c.Exists(key) {
		return nil, errors.New("key not present")
	}
	c.mu.RLock()
	data := c.Data[key]
	c.mu.RUnlock()
	if mpval, ok := data.Value.(map[string]interface{}); ok {
		if data, ok := mpval[field]; ok {
			return data, nil
		}
		return nil, errors.New("key not present")
	}
	return nil, errors.New("not a Hash value/table")

}

func (c *cache) HGetAll(key string) (map[string]interface{}, error) {
	c.mu.RLock()
	data, ok := c.Data[key]
	c.mu.RUnlock()
	if ok {

		if mpData, oks := data.Value.(map[string]interface{}); oks {
			return mpData, nil
		}
		return nil, errors.New("not a Hash value/table")
	}
	return nil, errors.New("key not present")
}

func (c *cache) DeleteExpiredKeys() {
	c.mu.RLock()
	hl := len(c.heap.Data)
	if hl == 0 {
		c.mu.RUnlock()
		return
	}
	node := c.heap.Data[0]
	c.mu.Unlock()
	if time.Now().Unix() > node.Expiry {
		c.mu.Lock()
		delete(c.Data, node.Key)
		c.heap.Extract()
		c.mu.Unlock()
	}
}
