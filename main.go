package main

import (
	"fmt"
	"gowi/gowi"
	"time"
)

func main() {
	c := gowi.NewCache()
	c.Set("name", 1000, "leoantony")

	tf := time.Second * 15
	time.Sleep(tf)
	// c.Del("name")

	// c.Hset("metadata", "age", 19)
	// c.Hset("metadata", "name", "leo")
	// c.Hset("metadata", "place", "ollur")

	// fmt.Println(c.HGetAll("metadata"))
	// fmt.Println(c.HGet("metadata","place"))

	fmt.Println(time.Now().Unix())
	val, _ := c.Get("name")

	fmt.Printf("%d \n", val)

	// if str,ok := val.(string); ok{
	// 	test(str)
	// }else{
	// 	fmt.Println("Not a string")
	// }

}


