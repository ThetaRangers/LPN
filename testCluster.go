package main

import (
	"SDCC/client"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

//test configuration const
const (
	stub          = false
	onClosestPeer = false

	//workers configuration
	workerPeriod  = 10
	minimumPeriod = 1

	workerKeySize       = 50
	workerThreads       = 50
	maxWorkerThreads    = 100
	increaseWorkerRatio = 10

	repetitions = 50

	putRatio = 50
	getRatio = 30
	//append ratio implicit

	//Client conf
	clientPeriod = 5

	//result files
	putFileName    = "put.csv"
	getFileName    = "get.csv"
	appendFileName = "append.csv"
)

var workerKeys []string

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

func trafficGeneratorThread(addr []string) {
	existingKey := workerKeys[0]
	for {
		addrToConnect := addr[rand.Intn(len(addr))]
		conn, grpcConn, _ := client.Connect(addrToConnect)

		time.Sleep(time.Duration(rand.Intn(workerPeriod))*time.Second + minimumPeriod*time.Second)

		x := rand.Intn(100)
		if x < putRatio {
			//putCase
			data := generateData(rand.Intn(100), rand.Intn(10)+1)
			key := workerKeys[rand.Intn(len(workerKeys))]
			err := conn.Put([]byte(key), data)
			if err != nil {
				log.Fatal(err)
			}
		} else if x < putRatio+getRatio {
			//getCase
			_, err := conn.Get([]byte(existingKey))
			if err != nil {
				log.Fatal(err)
			}
		} else {
			//appendCase
			data := generateData(rand.Intn(100), rand.Intn(10)+1)
			err := conn.Append([]byte(existingKey), data)
			if err != nil {
				log.Fatal(err)
			}
		}

		err := grpcConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func putMeasure(peers []string) float64 {

	var sum time.Duration
	sum = 0

	for i := 0; i < repetitions; i++ {
		time.Sleep(time.Duration(rand.Intn(clientPeriod))*time.Second + minimumPeriod*time.Second)
		addrToConnect := peers[rand.Intn(len(peers))]
		key := workerKeys[rand.Intn(len(workerKeys))]
		data := generateData(rand.Intn(100), 2)

		conn, grpcConn, _ := client.Connect(addrToConnect)
		start := time.Now()
		err := conn.Put([]byte(key), data)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		sum += elapsed

		err = grpcConn.Close()
		if err != nil {
			return 0
		}

	}
	return float64(sum.Milliseconds()) / float64(repetitions)
}

func getMeasure(peers []string) float64 {
	var sum time.Duration
	sum = 0

	for i := 0; i < repetitions; i++ {
		time.Sleep(time.Duration(rand.Intn(clientPeriod))*time.Second + minimumPeriod*time.Second)
		addrToConnect := peers[rand.Intn(len(peers))]
		key := workerKeys[rand.Intn(len(workerKeys))]
		conn, grpcConn, _ := client.Connect(addrToConnect)
		start := time.Now()
		_, err := conn.Get([]byte(key))
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		sum += elapsed

		err = grpcConn.Close()
		if err != nil {
			return 0
		}

	}
	return float64(sum.Milliseconds()) / float64(repetitions)
}

func appendMeasure(peers []string) float64 {
	var sum time.Duration
	sum = 0

	for i := 0; i < repetitions; i++ {
		time.Sleep(time.Duration(rand.Intn(clientPeriod))*time.Second + minimumPeriod*time.Second)
		addrToConnect := peers[rand.Intn(len(peers))]
		data := generateData(rand.Intn(100), 2)
		key := workerKeys[rand.Intn(len(workerKeys))]
		conn, grpcConn, _ := client.Connect(addrToConnect)
		start := time.Now()
		err := conn.Append([]byte(key), data)
		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		sum += elapsed

		err = grpcConn.Close()
		if err != nil {
			return 0
		}

	}
	return float64(sum.Milliseconds()) / float64(repetitions)
}

func main() {
	//addresses
	//needed to gateway or stubReg
	//TODO set confString
	confString := ""
	var networkAddr []string
	if stub {
		networkAddr, _ = client.GetAllNodesLocal(confString)
	} else {
		networkAddr = client.GetAllNodes(confString)
	}

	//workers key subset generation
	workerKeys = make([]string, workerKeySize)
	for i := 0; i < workerKeySize; i++ {
		workerKeys = append(workerKeys, strconv.FormatInt(int64(i), 10))
	}

	//nearest peer
	nearestPeer, err := client.GetClosestNode(networkAddr)
	if err != nil {
		log.Fatal(err)
	}

	//generate existing key
	conn, grpcConn, _ := client.Connect(nearestPeer)
	err = conn.Put([]byte(workerKeys[0]), generateData(rand.Intn(30), rand.Intn(5)+1))
	if err != nil {
		log.Fatal(err)
	}
	err = grpcConn.Close()
	if err != nil {
		log.Fatal(err)
	}

	//spawn thread
	for i := 0; i < workerThreads; i++ {
		go trafficGeneratorThread(networkAddr)
	}

	//set client peers
	var clientPeers []string
	if onClosestPeer {
		clientPeers = append(clientPeers, nearestPeer)
	} else {
		clientPeers = append(clientPeers, networkAddr...)
	}

	//setupFiles
	var headers = [][]string{{"Configuration"}, {"Network size: " + strconv.FormatInt(int64(len(networkAddr)), 10), "Using stub: " + strconv.FormatBool(stub)},
		{"Client requests on closes peer: " + strconv.FormatBool(onClosestPeer), "Worker period: " + strconv.FormatInt(workerPeriod, 10)},
		{"Client period: " + strconv.FormatInt(clientPeriod, 10), "Minimum period: " + strconv.FormatInt(minimumPeriod, 10)},
		{"Worker key size: " + strconv.FormatInt(workerKeySize, 10), "Max workers: " + strconv.FormatInt(maxWorkerThreads, 10)},
		{"Increase worker ratio: " + strconv.FormatInt(increaseWorkerRatio, 10), "Put ratio: " + strconv.FormatInt(putRatio, 10)},
		{"Get ratio: " + strconv.FormatInt(getRatio, 10), "Append ratio: " + strconv.FormatInt(100-getRatio-putRatio, 10)},
		{"Measure repetitions: " + strconv.FormatInt(repetitions, 10)}, {"\n"}, {"Operation", "Actual workers", "Response time (Seconds)"}}

	//create files and write headers
	putFile, err := os.OpenFile(putFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		log.Fatal(err)
	}
	defer func(putFile *os.File) {
		err := putFile.Close()
		if err != nil {

		}
	}(putFile)

	getFile, err := os.OpenFile(getFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		log.Fatal(err)
	}
	defer func(getFile *os.File) {
		err := getFile.Close()
		if err != nil {

		}
	}(getFile)

	appendFile, err := os.OpenFile(appendFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		log.Fatal(err)
	}
	defer func(appendFile *os.File) {
		err := appendFile.Close()
		if err != nil {

		}
	}(appendFile)

	putWriter := csv.NewWriter(putFile)
	defer putWriter.Flush()

	getWriter := csv.NewWriter(getFile)
	defer getWriter.Flush()

	appendWriter := csv.NewWriter(appendFile)
	defer appendWriter.Flush()

	for _, value := range headers {
		err := putWriter.Write(value)
		if err != nil {
			return
		}
		err = getWriter.Write(value)
		if err != nil {
			return
		}
		err = appendWriter.Write(value)
		if err != nil {
			return
		}
	}

	//run tests
	//TODO change format retValue
	for workerCount := workerThreads; workerCount < maxWorkerThreads; workerCount += increaseWorkerRatio {
		//Put test
		time.Sleep(time.Duration(rand.Intn(clientPeriod))*time.Second + minimumPeriod*time.Second)
		line := make([]string, 0)
		retValue := putMeasure(clientPeers)
		line = append(line, "Put", strconv.FormatInt(int64(workerCount), 10), fmt.Sprintf("%v", retValue))
		err := putWriter.Write(line)
		if err != nil {
			return
		}

		//Get test
		time.Sleep(time.Duration(rand.Intn(clientPeriod))*time.Second + minimumPeriod*time.Second)
		line = make([]string, 0)
		retValue = getMeasure(clientPeers)
		line = append(line, "Get", strconv.FormatInt(int64(workerCount), 10), fmt.Sprintf("%v", retValue))
		err = getWriter.Write(line)
		if err != nil {
			return
		}

		//Append test
		time.Sleep(time.Duration(rand.Intn(clientPeriod))*time.Second + minimumPeriod*time.Second)
		line = make([]string, 0)
		retValue = appendMeasure(clientPeers)
		line = append(line, "Append", strconv.FormatInt(int64(workerCount), 10), fmt.Sprintf("%v", retValue))
		err = appendWriter.Write(line)
		if err != nil {
			return
		}

		//increase workers
		for i := 0; i < increaseWorkerRatio; i++ {
			go trafficGeneratorThread(networkAddr)
			workerCount++
		}

	}
}
