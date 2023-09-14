package gowiz

// import (
// 	"fmt"
// 	"gowiz"
// 	"time"
// )

// func main() {
// 	c := gowiz.NewCache()
// 	c.Set("name", 2000, "leoantony")
// 	c.Set("jsondata", 10000, "THIS IS A TEST ")

// 	tf := time.Second * 25
// 	time.Sleep(tf)
// 	// c.Del("name")

// 	c.Hset("metadata", "age", 19)
// 	c.Hset("metadata", "name", "leo")
// 	c.Hset("metadata", "place", "ollur")

// 	fmt.Println(c.HGetAll("metadata"))
// 	fmt.Println(c.HGet("metadata", "place"))

// 	// fmt.Println(time.Now().Unix())
// 	val, _ := c.Get("name")
// 	val2, _ := c.Get("jsondata")

// 	fmt.Printf("%d \n", val)
// 	fmt.Printf("%d \n", val2)

// 	c.Set("test", 10000, "idk")

// 	c.Del("test")

// 	c.Set("hello", 10000, "idk2")

// 	c.Update("hello", "idk3")

// 	fmt.Println(c.Get("hello"))

// 	// if str,ok := val.(string); ok{
// 	// 	test(str)
// 	// }else{
// 	// 	fmt.Println("Not a string")
// 	// }

// 	// tg := make(chan int)

// 	// <-tg

// }
