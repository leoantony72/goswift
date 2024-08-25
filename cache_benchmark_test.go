package goswift
import (
	"testing"
)

func BenchmarkSet(b *testing.B){
	
	key := "name"
	val := "leoantony"
	cache := NewCache()
	b.ReportAllocs()
    for i:=0;i<b.N;i++{
		// for i:=0;i<100;i++{
			cache.Set(key,val,0)
		// }
	}
}
func BenchmarkSetWithExpiry(b *testing.B){
	
	key := "name"
	val := "leoantony"
	cache := NewCache()
	b.ReportAllocs()
    for i:=0;i<b.N;i++{
		// for i:=0;i<100;i++{
			cache.Set(key,val,10000)
		// }
	}
}
func BenchmarkGet(b *testing.B){
	
	key := "name"
	cache := NewCache()
	b.ReportAllocs()
    for i:=0;i<b.N;i++{
		// for i:=0;i<100;i++{
			cache.Get(key)
		// }
	}
}
// func BenchmarkArray(b *testing.B){
	
// 	key := []int{}
//     for i:=0;i<b.N;i++{
// 		for i:=0;i<100;i++{
// 			key=append(key,i)
// 		}
// 	}
// }