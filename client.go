package main

import (
	"SDCC/client"
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

func main() {

	allNodes, _ := client.GetAllNodesLocal("192.168.1.146")

	closest, _ := client.GetClosestNode(allNodes)

	conn, _, _ := client.Connect(closest)

	input := make([][]byte, 1)
	input[0] = []byte("defa")
	err := conn.Put([]byte("a"), input)
	if err != nil {
		log.Error(err)
	}

	time.Sleep(time.Second)

	val, err := conn.Get([]byte("a"))
	if err != nil {
		log.Error(err)
	}

	fmt.Println(val)

	err = conn.Append([]byte("a"), input)
	if err != nil {
		log.Error(err)
	}

	time.Sleep(time.Second)

	/*
		err = conn.Del([]byte("a"))
		if err != nil {
			log.Error(err)
		}*/

	val, err = conn.Get([]byte("a"))
	if err != nil {
		log.Error(err)
	}

	fmt.Println(val)
}
