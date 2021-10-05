package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"net"
	"net/http"
	"os"
)

var nodeId []byte
var targetId []byte

var m = 2

type metadataStruct struct {
	nodeId []byte
	ipAddress string
}

type stateTableStruct struct {
	leafSet [][]byte
	neighbourSet []metadataStruct
	routingTable [][]byte
}

var metadata metadataStruct
var stateTable stateTableStruct

func join(w http.ResponseWriter, req *http.Request) {
	if len(stateTable.neighbourSet) == 0 {
		//Nobody in the network
	}

	fmt.Fprintf(w, "%s\n", string(nodeId))
}

func connectToNode(arg string) {
	//Connect to a node
	fmt.Println("Connecting to localhost:", arg)

	targetIp := fmt.Sprintf("http://localhost:%s", arg)


	//NEED TO SEND JOIN(HASH)

	resp, err := http.Get(fmt.Sprintf("%s/join", targetIp))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		targetId = scanner.Bytes()
	}

	targetMetadata := metadataStruct{nodeId: targetId, ipAddress: targetIp}

	fmt.Println("Connected to: ", targetMetadata)

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

/*
Routing table
Leaf set: IP e ID secondo una metrica
Neighbour set: IP e ID di N nodi vicini
 */
func main() {
	//b := 4
	nodeId = make([]byte, 16)

	rand.Read(nodeId)
	fmt.Println("NodeId:", nodeId)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	metadata.nodeId = nodeId
	metadata.ipAddress = fmt.Sprintf("http://localhost:%s/nodeid", listener.Addr().(*net.TCPAddr).Port)

	//remove first one
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) != 0 {
		connectToNode(argsWithoutProg[0])
	}

	//Startup server
	http.HandleFunc("/join", join)

	fmt.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)

	http.Serve(listener, nil)
}
