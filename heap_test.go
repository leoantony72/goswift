package goswift

import (
	"fmt"
	"sync"
	"testing"

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
	cache.Set(key, 0, val)

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

func TestHeapExpiry(t *testing.T) {
	h := &expiry.Heap{}
	type th struct {
		key    string
		expiry int64
	}
	buildheap := []th{{"t1", 23}, {"t2", 12}, {"t3", 123}, {"t4", 436}, {"t5", 2}, {"t6", 14}, {"t7", 1}, {"t7", 6}}

	for _, v := range buildheap {
		h.Insert(v.key, v.expiry)
	}

	ExpectedValue := []int{1, 2, 6, 12, 14, 23, 123, 436}
	//Check heap values
	t.Run("HeapSortTest", func(t *testing.T) {
		for i := 0; i < len(ExpectedValue); i++ {
			val := Ex(h)
			if int(val.Expiry) != ExpectedValue[i] {
				t.Errorf("Expected Value: %d, Gotten : %d", ExpectedValue[i], val.Expiry)
				return
			}
			// fmt.Println(val.Expiry, ExpectedValue[i])
		}
	})
}

// func TestCache(t *testing.T) {
// 	c := goswift.NewCache()

// 	fmt.Println(time.Now().Unix())
// 	c.Set("leo", 23000, "kinglol")
// 	c.Set("name", 9000, "leoantony")
// 	c.Set("jsondata", 6000, "THIS IS A TEST ")
// 	exp := 3000
// 	var wg sync.WaitGroup
// 	for i := 0; i < 1000; i++ {
// 		wg.Add(3)
// 		go AddNode(c, exp, &wg)
// 		go AddNode(c, exp, &wg)
// 		go AddNode(c, exp, &wg)
// 	}
// 	c.Set("idk", 2000, "THIS IS A TEST ")
// 	c.Set("boiz", 7000, "THIS IS A TEST ")
// 	c.Set("no name", 10000, "THIS IS A TEST ")

// 	wg.Wait()

// 	PrintALL(c)

// 	c.Del("no name")
// 	interval := 1 * time.Second
// 	// fmt.Println(interval)
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			PrintALL(c)
// 			PrintALLH(c)
// 		}
// 	}
// }

func PrintALL(c CacheFunction) {
	d := c.AllData()
	fmt.Println(d)
	counTer := 0
	for _, v := range d {
		fmt.Println(v)
		counTer += 1
	}
	fmt.Println("total: ", counTer)
	fmt.Println("----------------------")
}
func PrintALLH(c CacheFunction) {
	d := c.AllDataHeap()
	// fmt.Println(d)
	counTer := 0
	for s, v := range d {
		fmt.Println(s, v)
		counTer += 1
	}
	fmt.Println("total Heap Data: ", counTer)
	fmt.Println("----------------------")
}

func AddNode(c CacheFunction, exp int, wg *sync.WaitGroup) {
	defer wg.Done()
	key := uuid.New()
	v := uuid.New()
	c.Set(key.String(), exp, v.String())
}

func Ex(h *expiry.Heap) *expiry.Node {
	c, _ := h.Extract()
	return c
}

func Print(h *expiry.Heap) {
	for _, b := range h.Data {
		fmt.Println(b)
	}
}
