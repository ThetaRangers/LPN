package main

import (
	"SDCC/go-dht/dht"
	"fmt"
)

func main() {

	client := dht.New(dht.DhtOptions{
		ListenAddr:    ":6000",
		BootstrapAddr: ":3000",
	})

	client.Start()

	hash, _, err := client.Store([]byte("Some value"))
	if err != nil {
		fmt.Println("ERROR1")
	}

	value, err := client.Fetch(hash)
	if err != nil {
		fmt.Println("ERROR2")
	}

	fmt.Println(string(value)) // Prints 'Some value'

	client.Stop()
}
