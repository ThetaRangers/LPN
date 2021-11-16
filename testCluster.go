package main

import (
	"SDCC/client"
	"fmt"
	"log"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateData(size int, rows int) [][]byte {
	rowSize := size / rows
	res := make([][]byte, 0)
	generate := make([]byte, rowSize)

	for i := 0; i < rows; i++ {
		rand.Read(generate)
		res = append(res, generate)
	}

	return res
}

func generateKey(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func workerThread(addr string) {
	conn, _, _ := client.Connect(addr)

	key := generateKey(10)
	data := generateData(30, 3)

	for {
		err := conn.Put([]byte(key), data)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

func main() {
	for i := 0; i < 50; i++ {
		go workerThread("172.17.0.2")
	}

	conn, _, _ := client.Connect("172.17.0.2")

	for {
		start := time.Now()
		err := conn.Put([]byte("a"), generateData(100, 2))
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)

		fmt.Println(elapsed)

		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	}
}
