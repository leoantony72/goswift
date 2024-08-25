package goswift

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	// "github.com/leoantony72/goswift"
	"github.com/leoantony72/goswift/expiry"
)

// const (
// 	ErrKeyNotFound  = "key does not Exists"
// 	ErrNotHashvalue = "not a Hash value/table"
// )

func TestSet(t *testing.T) {
	cache := NewCache()

	key := "name"
	val := "leoantony"
	cache.Set(key, val, 0)

	getValue, err := cache.Get(key)
	if err != nil {
		if err.Error() == ErrKeyNotFound {
			t.Errorf("key `%s`: %s", key, ErrKeyNotFound)
			return
		}
		return
	}
	if getValue.(string) != val {
		t.Errorf("val not the same")
		return
	}
}

func TestGet(t *testing.T) {
	c := NewCache()
	c.Set("age", 12, 0)

	val, err := c.Get("age")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if val.(int) != 12 {
		t.Errorf("Expected Value: 12(int) ,Gotten: %d", val)
		return
	}

	// Key does not exists
	_, err = c.Get("name")
	if err == nil {
		t.Errorf("Expected Error: %s", ErrKeyNotFound)
		return
	}

	//expiry provided- expiry>>time.Now()
	c.Set("place", "Kerala", 150000)
	val, err = c.Get("place")
	if err != nil {
		t.Errorf(err.Error())
		t.Errorf("key %s Might be expired", "place")
		return
	}
	if val.(string) != "Kerala" {
		t.Errorf("Expected Value: %s, Gotten: %s", "Kerala", val.(string))
		return
	}

	c.Set("country", "India", 100)
	_, err = c.Get("country")
	if err != nil {
		if err.Error() != ErrKeyNotFound {
			t.Errorf("key %s Should be Expired", "country")
			return
		}
	}

}

func TestUpdate(t *testing.T) {
	c := NewCache()

	key := "users:bob"
	value := "Cool shirt"
	c.Set(key, value, 0)

	data, err := c.Get(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if data.(string) != value {
		t.Errorf("Expected Value: %s ,Gotten: %s", value, data)
		return
	}

	newValue := "Chemistry sucks"
	c.Update(key, newValue)

	data, err = c.Get(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if data.(string) != newValue {
		t.Errorf("Expected Value: %s ,Gotten: %s", newValue, data)
		return
	}

	//key does not exist
	key = "water"
	err = c.Update(key, "H2O")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}

	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}
}

func TestDel(t *testing.T) {
	c := NewCache()
	key := "users:bob"
	value := "Cool shirt"
	c.Set(key, value, 0)

	ok := c.Exists(key)
	if !ok {
		t.Errorf("Expected Value: %v, Gotten: %v", true, ok)
		return
	}

	c.Del(key)
	ok = c.Exists(key)
	if ok {
		t.Errorf("Expected Value: %v, Gotten: %v", false, ok)
		return
	}

	//Key does not exist
	key = "users:varun"
	c.Del(key)

	// Key with Expiry
	key = "users:kingbob"
	c.Set(key, "bobbb!", 10000)
	c.Del(key)
}
func TestHset(t *testing.T) {
	c := NewCache()
	key := "users:John:metadata"
	c.Hset(key, "name", "John", 0)
	c.Hset(key, "age", 20, 0)
	c.Hset(key, "place", "Thrissur", 0)
	c.Hset(key, "people", []string{"bob", "tony", "henry"}, 0)

	data, err := c.HGetAll(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	name := data["name"].(string)
	age := data["age"].(int)
	place := data["place"].(string)
	people := data["people"].([]string)

	expectedArrayValues := []string{"bob", "tony", "henry"}

	if name != "John" {
		t.Errorf("Expected Value: %s, Gotten: %s", "John", name)
		return
	}

	if age != 20 {
		t.Errorf("Expected Value: %d, Gotten: %d", 20, age)
		return
	}

	if place != "Thrissur" {
		t.Errorf("Expected Value: %s, Gotten: %s", "Thrissur", place)
		return
	}

	i := 0
	t.Run("Hash :Array Data Type", func(t *testing.T) {

		for _, val := range expectedArrayValues {
			if val != people[i] {
				t.Errorf("Expected Value: %s, Gotten: %s", val, people[i])
				return
			}
			i++
		}
	})

	// hset with expiry non expired
	key = "users:user3"
	c.Hset(key, "test", "testvalue", 3000)

	_, err = c.HGet(key, "test")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	//hset with expiry expired key
	key = "users:user4"
	c.Hset(key, "test", "testvalue", 100)

	_, err = c.HGet(key, "test")
	if err == nil {
		t.Errorf("expected error: %s, Gotten: nil", ErrKeyNotFound)
		return
	}
	if err.Error() != ErrKeyNotFound {
		t.Errorf("expected error: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}

}

func TestHGet(t *testing.T) {
	c := NewCache()

	key := "users:Jhon:data"
	field := "age"
	value := 20
	c.Hset(key, field, value, 0)

	data, err := c.HGet(key, field)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if data.(int) != value {
		t.Errorf("Expected Value: %d, Gotten: %d", value, data)
		return
	}

	// key does not exists
	key = "fruits"
	_, err = c.HGet(key, "sweet")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}

	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}

	// field does not exist
	key = "fruits"
	field = "bitter"
	c.Hset(key, field, "lemons", 0)
	_, err = c.HGet(key, "sweet")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrFieldNotFound)
		return
	}

	if err.Error() != ErrFieldNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrFieldNotFound, err.Error())
		return
	}

	// Not an Hash Value
	key = "fruits"
	v := "orange"
	c.Set(key, v, 0)
	_, err = c.HGet(key, "sweet")
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrNotHashvalue)
		return
	}

	if err.Error() != ErrNotHashvalue {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrNotHashvalue, err.Error())
		return
	}
}

func TestHgetAll(t *testing.T) {
	opt := CacheOptions{
		EnableSnapshots:  true,
		SnapshotInterval: time.Second,
	}
	c := NewCache(opt)

	//key does not exists
	key := "users:bob"
	_, err := c.HGetAll(key)
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}

	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}

	//value not hash
	key = "fruits"
	v := "orange"
	c.Set(key, v, 0)
	_, err = c.HGetAll(key)
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrNotHashvalue)
		return
	}

	if err.Error() != ErrNotHashvalue {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrNotHashvalue, err.Error())
		return
	}

	//hgetall with expiry- non expired
	key = "users:user5"
	c.Hset(key, "test", "testvalue", 3000)

	_, err = c.HGetAll(key)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	//hgetall with expiry-  expired
	key = "users:user6"
	c.Hset(key, "test", "testvalue", 100)

	_, err = c.HGetAll(key)
	if err == nil {
		t.Errorf("Expected Err: %s, Gotten: ERR NIL", ErrKeyNotFound)
		return
	}
	if err.Error() != ErrKeyNotFound {
		t.Errorf("Expected Err: %s, Gotten: %s", ErrKeyNotFound, err.Error())
		return
	}

}

func TestHmset(t *testing.T) {
	c := NewCache()

	//Hmset struct test
	type ts struct {
		Name  string
		Age   int
		Place string
	}

	data := &ts{Name: "leo", Age: 23, Place: "ollur"}
	err := c.HMset("meta", data, 3000)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	// time.Sleep(time.Second * 4)
	ds, err := c.HGetAll("meta")
	if err != nil {
		fmt.Println(err)
		return
	}

	heapdata := testHeaps(c.(*Cache))
	if len(heapdata) > 1 {
		t.Errorf("expected heap length: 1 , gotten: %d", len(heapdata))
		return
	}
	if ds == nil {
		t.Errorf("expected to contain data,gotten nil")
		return
	}

	//hmset map test
	mapdata := map[string]interface{}{
		"name":  "loki",
		"place": "asgard",
		"age":   1054,
	}
	err = c.HMset("meta2", mapdata, 0)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	ds, err = c.HGetAll("meta2")
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	heapdata = testHeaps(c.(*Cache))
	if len(heapdata) > 1 {
		t.Errorf("expected heap length: 1 , gotten: %d", len(heapdata))
		return
	}
	if ds == nil {
		t.Errorf("expected to contain data,gotten nil")
		return
	}

	//invalid data type
	err = c.HMset("meta3", 34, 0)
	if err == nil {
		t.Errorf("expected to contain data,gotten nil")
		return
	}
	if err.Error() != ErrHmsetDataType {
		t.Errorf("expected error %s,gotten: %s", ErrHmsetDataType, err.Error())
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			id := uuid.New().String()
			c.HMset(id, mapdata, 0)
			wg.Done()
		}()
		go func() {
			id := uuid.New().String()
			c.HMset(id, mapdata, 0)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestExist(t *testing.T) {
	c := NewCache()

	key := "users:bob"
	c.Set(key, "mexican alien", 4000)

	ok := c.Exists(key)
	if !ok {
		t.Errorf("Expected Value: %v, Gotten: %v", true, ok)
		return
	}

	key = "users:john"
	t.Run("Key does not exist", func(t *testing.T) {
		ok := c.Exists(key)
		if ok {
			t.Errorf("Expected Value: %v, Gotten: %v", false, ok)
			return
		}
	})
}

func TestGetAllData(t *testing.T) {
	c := NewCache()

	keys := []string{"name", "age", "idk"}
	c.Set(keys[0], "bob", 0)
	c.Set(keys[1], 22, 0)
	c.Set(keys[2], "idk", 0)
	data, _ := c.AllData()

	for i := 0; i < len(keys); i++ {
		if _, ok := data[keys[i]]; !ok {
			t.Errorf("Key:%s does't Exist", keys[i])
		}
	}

}

func TestDeleteExpiredKeys(t *testing.T) {
	c := NewCache().(*Cache)

	c.Set("key1", "t1", 100)
	c.Set("key2", "t1", 200)
	c.Set("key3", "t1", 3200)

	time.Sleep(time.Second * 1)
	testDeleteExpiredKeys(c)

	if c.Exists("key1") || c.Exists("key2") {
		t.Errorf("key1 & key2 has not been removed")
		return
	}
	if !c.Exists("key3") {
		t.Errorf("key3 should exists")
	}
	c.Del("key3")
	testDeleteExpiredKeys(c)

	c.Set("key4", "t4", 0)
	testDeleteExpiredKeys(c)

}

func AddNode(c CacheFunction, exp int, wg *sync.WaitGroup) {
	defer wg.Done()
	key := uuid.New()
	v := uuid.New()
	c.Set(key.String(), v.String(), exp)
}

func Print(h *expiry.Heap) {
	for _, b := range h.Data {
		fmt.Println(b)
	}
}

func TestSnapshotWithoutOpt(t *testing.T) {
	// opt := CacheOptions{
	// 	EnableSnapshots:  true,
	// 	SnapshotInterval: time.Second,
	// }
	c := NewCache()

	c.Set("user:1", "bob", 0)
	c.Hset("user:2", "name", "jhon", 0)
	c.Set("user:3", "raju", 3000000)

	snapshot(c.(*Cache))
	c.Del("user:1")
	c.Del("user:2")
	c.Del("user:3")
	fmt.Println(c.AllData())

	decoder(c.(*Cache))

	fmt.Println(c.AllData())
	if !c.Exists("user:1") || !c.Exists("user:2") || !c.Exists("user:3") {
		t.Errorf("Key does not exists, Snapshot does not take place")
		return
	}

	// time.Sleep(time.Second * 5)
}

func TestSnapshotWithOpt(t *testing.T) {
	opt := CacheOptions{
		EnableSnapshots:  true,
		SnapshotInterval: time.Second,
	}
	c := NewCache(opt)
	c.Set("user:1", "bob", 0)
	c.Hset("user:2", "name", "jhon", 0)
	c.Set("user:3", "raju", 3000000)
}

func TestSnapshotTimer(t *testing.T) {
	c := NewCache()

	Close := make(chan struct{})
	go snapShotTimer(c.(*Cache), time.Millisecond,Close)
	close(Close)
}

func TestDecoder(t *testing.T) {
	opt := CacheOptions{
		EnableSnapshots:  true,
		SnapshotInterval: time.Millisecond,
	}
	c := NewCache(opt)

	c.Set("u1","lol",0)
	c.Set("u2","lol",0)
	c.Set("u3","lol",0)
	c.Set("u4","lol",0)
	time.Sleep(time.Millisecond*1000)
	testdecoder()
}