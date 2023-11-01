package goswift

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

func SnapShotTimer(c *Cache, t time.Duration) {
	_, err := os.Create("snapshot.data")
	if err != nil {
		log.Fatal(err)
		return
	}
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			Snapshot(c)
		}
	}

}

func Snapshot(c *Cache) {
	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)
	data := c.AllDatawithExpiry()
	fmt.Println("MapData: ", data)

	if err := enc.Encode(data); err != nil {
		log.Fatal("err snapshot: ", err)
	}

	// file, err := os.OpenFile("snapshot.data", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file, err := os.Create("snapshot.data")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	n, err := file.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("count:", n)
}

func Decoder(c *Cache) {
	file, err := os.Open("snapshot.data")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data := make(map[string]SnaphotData)
	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		fmt.Println(err)
	}
}
