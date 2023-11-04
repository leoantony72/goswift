package goswift

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

func SnapShotTimer(c *Cache, t time.Duration) {
	_, err := os.Create("snapshot.data")
	if err != nil {
		// log.Fatal(err)
		return
	}
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			Snapshot(c)
		case <-Close:
			break
		}
	}

}

func Snapshot(c *Cache) {
	var buffer bytes.Buffer

	gob.Register(map[string]interface{}{})
	enc := gob.NewEncoder(&buffer)
	data := c.AllDatawithExpiry()
	// fmt.Println("MapData: ", data)

	if err := enc.Encode(data); err != nil {
		fmt.Println("err snapshot: ", err)
		return
	}

	file, err := os.Create("snapshot.data")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println("count:", n)
}

func Decoder(c *Cache) {
	gob.Register(map[string]interface{}{})
	file, err := os.Open("snapshot.data")
	if err != nil {
		fmt.Println("file open err: ", err)
		// os.Create("snapshot.data")
		// file, _ = os.Open("snapshot.data")
		return
	}
	defer file.Close()

	data := make(map[string]SnaphotData)
	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		fmt.Println("decode err", err)
	}
	fmt.Println("decoded data", data)
	AddToCache(data, c)
}

func AddToCache(d map[string]SnaphotData, c *Cache) {
	for k, v := range d {
		c.Set(k, v.Value, int(v.Expiry))
	}
	// fmt.Println("cache Data: ", c.AllDatawithExpiry())
}
