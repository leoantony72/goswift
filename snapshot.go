package goswift

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

func snapShotTimer(c *Cache, t time.Duration, fileName string, Close chan struct{}) {
	_, err := os.Create(fileName)
	if err != nil {
		// log.Fatal(err)
		return
	}
	ticker := time.NewTicker(t)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			snapshot(c, fileName)
		case <-Close:
			return
		}
	}

}

func snapshot(c *Cache, fileName string) {
	var buffer bytes.Buffer

	gob.Register(map[string]interface{}{})
	enc := gob.NewEncoder(&buffer)
	data := c.AllDatawithExpiry()

	if err := enc.Encode(data); err != nil {
		fmt.Println("err snapshot: ", err)
		return
	}

	file, err := os.Create(fileName)
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
}
func decoder(c *Cache, fileName string) {
	gob.Register(map[string]interface{}{})
	file, err := os.Open("snapshot.data")
	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create("snapshot.data")
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
		} else {
			fmt.Println("Error opening file:", err)
			return
		}
	}
	defer file.Close()

	// Check if the file is empty
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}
	if fileInfo.Size() == 0 {
		fmt.Println("File is empty, nothing to decode")
		return
	}

	data := make(map[string]snapShotData)
	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		fmt.Println("Decode error:", err)
		return
	}

	addToCache(data, c)
}

func testdecoder() {
	gob.Register(map[string]interface{}{})
	file, err := os.Open("snapshot.data")
	if err != nil {
		os.Create("snapshot.data")
		return
	}
	defer file.Close()

	// Check if the file is empty before decoding
	fileInfo, err := file.Stat()
	if err != nil {
		// fmt.Println("file stat error:", err)
		return
	}
	if fileInfo.Size() == 0 {
		// fmt.Println("File is empty. Nothing to decode.")
		return
	}

	data := make(map[string]snapShotData)
	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&data); err != nil {
		fmt.Println("decode err", err)
	}
	// addToCache(data, c)
	fmt.Println(data)
}
func addToCache(d map[string]snapShotData, c *Cache) {
	for k, v := range d {
		c.Set(k, v.Value, int(v.Expiry))
	}
	// fmt.Println("cache Data: ", c.AllDatawithExpiry())
}
