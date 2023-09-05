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

func (c *Cache) Set(key string, val interface{}) {
	c.data[key] = val
}

func (c *Cache) Del(key string) {
	delete(c.data, key)
}

func (c *Cache) Update(key string, val interface{}) error {
	//check if key is present
	if _, ok := c.data[key]; ok {
		return errors.New("key not present")
	}
	c.data[key] = val
	return nil
}
