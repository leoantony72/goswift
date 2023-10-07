package goswift

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/leoantony72/goswift/expiry"
)

const (
	ErrKeyNotFound   = "key does not Exists"
	ErrFieldNotFound = "field does not Exists"
	ErrNotHashvalue  = "not a Hash value/table"
	ErrHmsetDataType = "invalid data type, Expected Struct/Map"
)

type Cache struct {
	Data   map[string]*dataHolder
	length int
	heap   *expiry.Heap
	mu     sync.RWMutex
}

type dataHolder struct {
	Value  interface{}
	Expiry *expiry.Node
}

// func (c *Cache) AllDataHeap() []*expiry.Node {
// 	c.mu.Lock()
// 	// var h []*expiry.Heap
// 	// d := c.heap.Data
// 	dst := make([]*expiry.Node, len(c.heap.Data))

// 	copy(dst, c.heap.Data)
// 	c.mu.Unlock()
// 	return dst
// }

// returns all data from the map with both key and value,
// expiry data will not be returned, returned data will be a
// copy of the original data
func (c *Cache) AllData() (map[string]interface{}, int) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dataMap := make(map[string]interface{}, len(c.Data))
	counter := 0
	for k, v := range c.Data {
		dataMap[k] = v.Value
		counter++
	}

	return dataMap, counter
}

// Initialize a New CacheFunction type which is an Interfaces
// for all the availabel function
func NewCache() CacheFunction {
	dataMap := make(map[string]*dataHolder)
	heapInit := expiry.Init()
	cache := &Cache{Data: dataMap, length: 0, heap: heapInit}
	go sweaper(cache, heapInit)
	return cache
}

func testHeaps(c *Cache) []*expiry.Node {
	c.mu.Lock()
	// var h []*expiry.Heap
	// d := c.heap.Data
	dst := make([]*expiry.Node, len(c.heap.Data))

	copy(dst, c.heap.Data)
	c.mu.Unlock()
	return dst
}

// Exists func receives the key check if it exists in the
// Hash Table and returns a boolean
func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	_, ok := c.Data[key]
	c.mu.RUnlock()
	return ok
}

// Internal exist func wihtout locking(c.mu.Lock)
func (c *Cache) ExistsNonBlocking(key string) bool {
	_, ok := c.Data[key]
	return ok
}

// Adds an element to Hash Set, If exp is provided add the
// Node to the Heap with Key and expiration time(int64).
// If exp == 0, Item Never expires, thus it isn't added
// In the Heap
func (c *Cache) Set(key string, val interface{}, exp int) {
	c.mu.Lock()
	var node *expiry.Node
	if exp == 0 {
		data := &dataHolder{Value: val}
		c.Data[key] = data
		c.mu.Unlock()
		return
	}
	exp = exp / 1000
	expTime := time.Now().Add(time.Second * time.Duration(exp)).Unix()
	node = c.heap.Insert(key, expTime)
	data := &dataHolder{Value: val, Expiry: node}
	c.Data[key] = data
	c.mu.Unlock()
}

// If exp is not nil, check if the element has expired or not
// Removes the element from the cache if expired, It does not remove
// The Node from the Heap, which will handled by the Sweaper.
// @This must be improved, So that deleted keys does'nt stay in the Heap.
func (c *Cache) Get(key string) (interface{}, error) {
	c.mu.RLock()
	data, ok := c.Data[key]
	c.mu.RUnlock()
	if !ok {
		return nil, errors.New(ErrKeyNotFound)
	}
	if data.Expiry != nil {
		if data.Expiry.Expiry > time.Now().Unix() {
			return data.Value, nil
		}
		c.mu.Lock()
		delete(c.Data, key)
		c.mu.Unlock()
		return nil, errors.New(ErrKeyNotFound)
	}
	return data.Value, nil
}

// Delete's the Item from the cache
// @This must be improved, So that deleted keys does'nt stay in the Heap.
func (c *Cache) Del(key string) {
	c.mu.Lock()
	if !c.ExistsNonBlocking(key) {
		c.mu.Unlock()
		return
	}
	data := c.Data[key]
	// if !ok {
	// 	c.mu.Unlock()
	// 	return
	// }

	fmt.Println(data.Expiry)
	if data.Expiry == nil {
		delete(c.Data, key)
		c.mu.Unlock()
		return
	}

	index := data.Expiry.Index
	c.heap.Remove(index, len(c.heap.Data)-1)
	delete(c.Data, key)
	c.mu.Unlock()
}

// Set a new value for the key only if it already exist.
// New data will expire at the same time as the prev Key.
func (c *Cache) Update(key string, val interface{}) error {
	if !c.Exists(key) {
		return errors.New(ErrKeyNotFound)
	}

	c.mu.Lock()
	e := c.Data[key].Expiry
	data := &dataHolder{Value: val, Expiry: e}
	c.Data[key] = data
	c.mu.Unlock()
	return nil
}

// Support for Hash data type, Hset func receives key, field and value
func (c *Cache) Hset(key, field string, value interface{}, exp int) {

	c.mu.Lock()
	if _, exists := c.Data[key]; !exists {
		c.Data[key] = &dataHolder{}
		c.Data[key].Value = make(map[string]interface{})
	} else {
		if c.Data[key].Expiry != nil {
			c.heap.Remove(c.Data[key].Expiry.Index, len(c.heap.Data)-1)
		}
	}
	var node *expiry.Node

	if exp != 0 {
		exp = exp / 1000
		expTime := time.Now().Add(time.Second * time.Duration(exp)).Unix()
		node = c.heap.Insert(key, expTime)
		c.Data[key].Expiry = node
	}

	hash := c.Data[key].Value.(map[string]interface{})
	hash[field] = value
	c.mu.Unlock()
}

// Retrieves the field value of hash by key and field name
func (c *Cache) HGet(key, field string) (interface{}, error) {
	if !c.Exists(key) {
		return nil, errors.New(ErrKeyNotFound)
	}
	c.mu.RLock()
	data := c.Data[key]
	c.mu.RUnlock()
	if mpval, ok := data.Value.(map[string]interface{}); ok {

		if dataf, ok := mpval[field]; ok {
			if data.Expiry != nil {
				if data.Expiry.Expiry > time.Now().Unix() {
					return dataf, nil
				}
				c.mu.Lock()
				delete(c.Data, key)
				c.mu.Unlock()
				return nil, errors.New(ErrKeyNotFound)
			}
			return dataf, nil
		}
		return nil, errors.New(ErrFieldNotFound)
	}
	return nil, errors.New(ErrNotHashvalue)

}

/* @struct
type t struct{
	name string
	age int
	place string
}

type tm map[string]int

tm{}
*/

// HMset takes a Struct/Map as the value , if any other datatype is 
// Provided it will return an error.
func (c *Cache) HMset(key string, d interface{}, exp int) error {
	valType := reflect.TypeOf(d)
	fieldValues := reflect.ValueOf(d)

	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
		fieldValues = fieldValues.Elem()
	}
	fmt.Println("type", valType.Kind())
	switch valType.Kind() {
	case reflect.Struct:
		{
			c.Data[key] = &dataHolder{Value: make(map[string]interface{})}
			for i := 0; i < valType.NumField(); i++ {
				field := valType.Field(i)
				value := fieldValues.Field(i).Interface()
				c.Hset(key, field.Name, value, exp)
			}
		}
	case reflect.Map:
		{
			for field, value := range d.(map[string]interface{}) {
				c.Hset(key, field, value, exp)
			}
		}
	default:
		return errors.New(ErrHmsetDataType)
	}
	return nil
}

// HgetAll retrives all the fields in a Hash by providing the key
func (c *Cache) HGetAll(key string) (map[string]interface{}, error) {
	c.mu.RLock()
	data, ok := c.Data[key]
	c.mu.RUnlock()
	if ok {

		if mpData, oks := data.Value.(map[string]interface{}); oks {
			if data.Expiry != nil {
				if data.Expiry.Expiry > time.Now().Unix() {
					return mpData, nil
				}
				c.mu.Lock()
				delete(c.Data, key)
				c.mu.Unlock()
				return nil, errors.New(ErrKeyNotFound)
			}
			return mpData, nil
		}
		return nil, errors.New(ErrNotHashvalue)
	}
	return nil, errors.New(ErrKeyNotFound)
}
