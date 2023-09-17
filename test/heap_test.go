package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/leoantony72/goswift"
	"github.com/leoantony72/goswift/expiry"
)

func TestHeap(t *testing.T) {
	h := &expiry.Heap{}

	type th struct {
		key    string
		expiry int64
	}
	buildheap := []th{{"t1", 23}, {"t2", 12}, {"t3", 123}, {"t4", 436}, {"t5", 2}, {"t6", 14}, {"t7", 1}}

	for _, v := range buildheap {
		h.Insert(v.key, v.expiry)
	}

	// for _,b := range h.Data{
	// 	fmt.Println(b)
	// }
	// t.Log(h)
	Ex(h)
	// fmt.Println("VALUE: ", c)
	Print(h)
	fmt.Println("-----------")

	Ex(h)
	// fmt.Println("VALUE: ", c)
	Print(h)
	fmt.Println("-----------")
	// t.Log(h.Extract())
	// fmt.Println(h)
	// t.Log(h)

}

func TestCache(t *testing.T) {
	c := goswift.NewCache()

	fmt.Println(time.Now().Unix())
	c.Set("leo", 23000, "kinglol")
	c.Set("name", 9000, "leoantony")
	c.Set("jsondata", 6000, "THIS IS A TEST ")
	exp := 3000
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(3)
		go AddNode(c, exp, &wg)
		go AddNode(c, exp, &wg)
		go AddNode(c, exp, &wg)
	}
	c.Set("idk", 2000, "THIS IS A TEST ")
	c.Set("boiz", 7000, "THIS IS A TEST ")
	c.Set("no name", 10000, "THIS IS A TEST ")

	wg.Wait()

	PrintALL(c)

	// time.Sleep(time.Second * 20)
	// fmt.Println(c.AllData())
	// PrintALL(c)

	// Print()

	c.Del("no name")
	interval := 1 * time.Second
	// fmt.Println(interval)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			PrintALL(c)
			PrintALLH(c)
		}
	}
}

func PrintALL(c goswift.CacheFunction) {
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
func PrintALLH(c goswift.CacheFunction) {
	d := c.AllDataHeap()
	// fmt.Println(d)
	counTer := 0
	for s, v := range d {
		fmt.Println(s,v)
		counTer += 1
	}
	fmt.Println("total Heap Data: ", counTer)
	fmt.Println("----------------------")
}

func AddNode(c goswift.CacheFunction, exp int, wg *sync.WaitGroup) {
	defer wg.Done()
	key := uuid.New()
	v := uuid.New()
	c.Set(key.String(), exp, v.String())
}

func Ex(h *expiry.Heap) {
	c, _ := h.Extract()
	fmt.Println("VALUE: ", c)
}

func Print(h *expiry.Heap) {
	for _, b := range h.Data {
		fmt.Println(b)
	}
}
