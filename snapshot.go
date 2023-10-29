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
	fmt.Println(data)

	if err := enc.Encode(data); err != nil {
		log.Fatal("errr snapshot: ", err)
	}

	file, err := os.OpenFile("snapshot.data", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	fmt.Println(n)

}