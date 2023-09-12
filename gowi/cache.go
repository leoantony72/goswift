package gowi

import (
	"errors"
	"fmt"
	"gowi/gowi/expiry"
	"time"
)

type Cache struct {
	data   map[string]*DataHolder
	length int
	heap   *expiry.Heap
}

type DataHolder struct {
	val    interface{}
	expiry *expiry.Node
}

func NewCache() *Cache {
	datamap := make(map[string]*DataHolder)
	heap := expiry.Init()
	return &Cache{data: datamap, length: 0, heap: heap}
}

func (c *Cache) Exists(key string) bool {
	if _, ok := c.data[key]; ok {
		return true
	}
	return false
}

func (c *Cache) Set(key string, exp int, val interface{}) {
	var node *expiry.Node
	if exp != 0 {
		exp = exp / 1000
		expTime := time.Now().Add(time.Second * time.Duration(exp)).Unix()
		node = c.heap.Insert(expTime, key)
		datamap := &DataHolder{val: val, expiry: node}
		c.data[key] = datamap
		return
	}
	datamap := &DataHolder{val: val}
	c.data[key] = datamap
	fmt.Println(c.data[key])
}

func (c *Cache) Get(key string) (interface{}, error) {
	val, ok := c.data[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}

	if val.expiry != nil {
		fmt.Println(val.expiry.Expiry, time.Now().Unix())
		if val.expiry.Expiry > time.Now().Unix() {
			return val, nil
		}
		delete(c.data, key)
		return nil, errors.New("key does not exist")
	}
	return val.val, nil
}

func (c *Cache) Del(key string) {
	delete(c.data, key)
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
	c.data[key] = datamap
	return nil
}

func (c *Cache) Hset(key, field string, value interface{}) {

	if _, exists := c.data[key]; !exists {
		c.data[key] = &DataHolder{}
		c.data[key].val = make(map[string]interface{})
	}

	hash := c.data[key].val.(map[string]interface{})
	hash[field] = value
}

func (c *Cache) HGet(key, field string) (interface{}, error) {
	if !c.Exists(key) {
		return nil, errors.New("key not present")
	}
	val, _ := c.data[key]

	if mpval, ok := val.val.(map[string]interface{}); ok {
		if data, ok := mpval[field]; ok {
			return data, nil
		}
		return nil, errors.New("key not present")
	}
	return nil, errors.New("not a Hash value/table")

}

func (c *Cache) HGetAll(key string) (map[string]interface{}, error) {
	if data, ok := c.data[key]; ok {

		if mpdata, oks := data.val.(map[string]interface{}); oks {
			return mpdata, nil
		}
		return nil, errors.New("not a Hash value/table")
	}
	return nil, errors.New("key not present")
}
