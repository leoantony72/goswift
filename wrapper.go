package goswift

type CacheFunction interface {
	Exists(key string) bool
	Set(key string, val interface{}, exp int)
	Get(key string) (interface{}, error)
	Del(key string)
	Update(key string, val interface{}) error
	Hset(key, field string, value interface{}, exp int)
	HGet(key, field string) (interface{}, error)
	HGetAll(key string) (map[string]interface{}, error)
	HMset(key string, d interface{}, exp int) error
	AllData() (map[string]interface{}, int)
	// AllDataHeap() []*expiry.Node
}
