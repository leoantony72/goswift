# GoSwift - Embedded cache for golang

High-performance, concurrent embedded caching library for Go applications with support for Hash data type

## Features

- Set & Get command
- Del command
- Update command
- Exists command
- Support for TTL
- Support for Disk Save(Snapshots)
- Support Hash Data type Hset, Hget, HgetAll, HMset
- Safe Locking

## Installation

```shell
go mod init github.com/my/repo
```

Then install goswift:

```shell
go get github.com/leoantony72/goswift
```

## Quickstart

```go

package main

import (
    "fmt"
    "github.com/leoantony72/goswift"
)

func main(){
    cache := goswift.NewCache()

    // Value 0 indicates no expiry
    cache.Set("key", "value", 0)

    val, err := cache.Get("key")
    if err !=nil{
        fmt.Println(err)
        return
    }
    fmt.Println("key", val)
}

```

## Disk Save

### Snapshot

```go
opt := goswift.CacheOptions{
		EnableSnapshots:  true,
		SnapshotInterval: time.Second*5,
	}
c := goswift.NewCache(opt)
```
> **_NOTE:_** If the **EnableSnapshot** is **false**, Data saved in the file will not imported

This will take a snapshot of the Data Every 5sec and saves it into a ***Snapshot.data*** file. By default Snapshots are disabled and if the SnapshotInterval is not provided default value is **5seconds**.

> **_NOTE:_** Don't delete the ***Snapshot.data*** File <br>

## Error Handling
```go
const (
	ErrKeyNotFound   = "key does not Exists"
	ErrFieldNotFound = "field does not Exists"
	ErrNotHashvalue  = "not a Hash value/table"
	ErrHmsetDataType = "invalid data type, Expected Struct/Map"
)
```
These are the common Errors that may occur while writing the code. These Varible provide you a clear and easy **Error** comparison method to determine errors.

```go
data,err := cache.Get("key")
if err != nil {
	if err.Error() == goswift.ErrKeyNotFound {
        //do something
}
}    
```
## Usage

```go
// Set Value with Expiry
// @Set(key string, val interface{}, exp int)
// Here expiry is set to 1sec
cache.Set("key","value",1000)


// Get Value with key
// @Get(key string) (interface{}, error)
val,err := cache.Get("key")
if err != nil{
    fmt.Println(err)
    return
}


// Update value
// @Update(key string, val interface{}) error
err = cache.Update("key","value2")
if err != nil{
    fmt.Println(err)
    return
}


// Delete command
// @Del(key string)
cache.Del("key")


// Hset command
// @Hset(key, field string, value interface{}, exp int)
// in this case the "key" expires in 1sec
cache.Hset("key","name","value",1000)
cache.Hset("key","age",18,1000)


// HMset command
// @HMset(key string, d interface{}, exp int) error
// Set a Hash by passing a Struct/Map
// ---by passing a struct---
type Person struct{
    Name  string
    Age   int
    Place string
}

person1 := &Person{Name:"bob",Age:18,Place:"NYC"}
err = cache.HMset("key",person1)
if err != nil{
    fmt.Println(err)
    return
}

// ---by passing a map---
person2 := map[string]interface{Name:"john",Age:18,Place:"NYC"}
err = cache.HMset("key",person2)
if err != nil{
    fmt.Println(err)
    return
}


// Hget command
// @HGet(key, field string) (interface{}, error)
// get individual fields in Hash
data,err := cache.HGet("key","field")
if err != nil{
    fmt.Println(err)
    return
}
fmt.Println(data)

// HgetAll command
// @HGetAll(key string) (map[string]interface{}, error)
// gets all the fields with value in a hash key
// retuns a map[string]interface{}
data,err = cache.HGetAll("key")
if err != nil{
    fmt.Println(err)
    return
}


// Exist command
// @Exists(key string) bool
// Check if the key exists
value = cache.Exists("key")
fmt.Println(value)



// AllData command
// @AllData() (map[string]interface{}, int)
// returns all the data in the cache with keys, also with no.of keys present
// returns the value as a map[strirng]interface{}
// !It does not return the expiry time of the key
data,counter := cache.AllData()
fmt.Println(data,counter)

```

## Run the Test

```go
go test ./...
```
