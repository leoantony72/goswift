package test

import (
	"fmt"
	"github.com/leoantony72/goswift"
	"github.com/leoantony72/goswift/expiry"
	"testing"
	"time"
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
	c.Set("leo", 3000, "kinglol")
	c.Set("name", 2000, "leoantony")
	c.Set("jsondata", 10000, "THIS IS A TEST ")

	// time.Sleep(time.Second*14)

	// Print(c)

	ts := make(chan int)
	ts <- -1
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
