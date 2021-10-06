package main

import (
	"bufio"
	"crypto/sha1"
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

	//fmt.Fprintf(w, "%s\n", string(nodeId))
	fmt.Fprintf(w, "%s\n",string(metadata.nodeId))
}

func connectToNode(fullAddress string) {
	//Connect to a node
	fmt.Println("Connecting to : ", fullAddress)

	targetIp := fmt.Sprintf(fullAddress)


	//NEED TO SEND JOIN(HASH)

	resp, err := http.Get(fmt.Sprintf("http://%s/join", targetIp))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		targetId = []byte(scanner.Text())
	}

	targetMetadata := metadataStruct{nodeId: targetId, ipAddress: targetIp}

	fmt.Printf("Connected to fullAddress %s nodeId %x\n", fullAddress, targetMetadata.nodeId)

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

/*
Routing table
Leaf set: IP e ID secondo una metrica
Neighbour set: IP e ID di N nodi vicini
 */

/****************************
arg[0] -> IP
arg[1] -> port
arg[2] -> debug connect/server
*****************************/
func main() {
	//b := 4
	nodeId = make([]byte, 16)

	//remove first one
	argsWithoutProg := os.Args[1:]

	ip := argsWithoutProg[0]
	port := argsWithoutProg[1]

	fullAddress := fmt.Sprintf("%s:%s", ip, port)

	//nodeId with hash
	h := sha1.New()
	h.Write([]byte(fullAddress))
	nodeId := h.Sum(nil)

	fmt.Printf("fullAddress %s nodeId %x\n", fullAddress, nodeId)

	listener, err := net.Listen("tcp", fullAddress)
	if err != nil {
		panic(err)
	}

	metadata.nodeId = nodeId
	metadata.ipAddress = fullAddress


	if len(argsWithoutProg) == 3 {
		connectToNode(argsWithoutProg[2])
	}

	//Startup server
	http.HandleFunc("/join", join)

	http.Serve(listener, nil)
}
