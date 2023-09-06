package gowi

import "errors"

type Cache struct {
	data   map[string]interface{}
	length int
}

func NewCache() *Cache {
	datamap := make(map[string]interface{})
	return &Cache{data: datamap, length: 0}
}

func (c *Cache) Exists(key string) bool {
	if _, ok := c.data[key]; ok {
		return true
	}
	return false
}

func (c *Cache) Set(key string, val interface{}) {
	c.data[key] = val
}

func (c *Cache) Get(key string) (interface{}, error) {
	val, ok := c.data[key]
	if !ok {
		return nil, errors.New("key does not exist")
	}
	return val, nil
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
	c.data[key] = val
	return nil
}

func (c *Cache) Hset(key, field string, val interface{}) {

	if _, exists := c.data[key]; !exists {
		c.data[key] = make(map[string]interface{})
	}

	hash := c.data[key].(map[string]interface{})
	hash[field] = val
}

func (c *Cache) HGet(key, field string) (interface{}, error) {
	if !c.Exists(key) {
		return nil, errors.New("key not present")
	}
	val, _ := c.data[key]

	if mpval, ok := val.(map[string]interface{}); ok {
		if data, ok := mpval[field]; ok {
			return data, nil
		}
		return nil, errors.New("key not present")
	}
	return nil, errors.New("not a Hash value/table")

}

func (c *Cache) HGetAll(key string) (map[string]interface{}, error) {
	if data, ok := c.data[key]; ok {

		if mpdata, oks := data.(map[string]interface{}); oks {
			return mpdata, nil
		}
		return nil, errors.New("not a Hash value/table")
	}
	return nil, errors.New("key not present")
}
