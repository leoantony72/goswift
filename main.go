package main

import (
	"fmt"
	"gowi/gowi"
)

func main() {
	c := gowi.NewCache()
	// c.Set("name", 213)

	// c.Del("name")

	c.Hset("metadata", "age", 19)
	c.Hset("metadata", "name", "leo")
	c.Hset("metadata", "place", "ollur")

	fmt.Println(c.HGetAll("metadata"))
	fmt.Println(c.HGet("metadata","place"))

	// val, _ := c.Get("name")

	// fmt.Printf("%T \n", val)

	// if str,ok := val.(string); ok{
	// 	test(str)
	// }else{
	// 	fmt.Println("Not a string")
	// }

}

// func test(v string) {
// 	fmt.Println(v)
// }
