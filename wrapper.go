package goswift

type CacheFunction interface {
	Exists(key string) bool
	Set(key string, exp int, val interface{})
	Get(key string) (interface{}, error)
	Del(key string)
	Update(key string, val interface{}) error
	Hset(key, field string, value interface{})
	HGet(key, field string) (interface{}, error)
	HGetAll(key string) (map[string]interface{}, error)
	AllData() (map[string]interface{}, int)
	// AllDataHeap() []*expiry.Node
}
