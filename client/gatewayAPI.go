package main

import (
	op "SDCC/operations"
	pb "SDCC/registerServer"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	//TODO set as terraform output result
	url = ""
)

type Connection struct {
	c op.OperationsClient
}

type IpListStruct struct {
	IpList []string `json:"ipList"`
}

func GetAllNodes() []string {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	if err != nil {
		log.Error(err)
	}

	//Close Body after read
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(response.Body)

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
	}

	ipList := IpListStruct{}

	if err := json.Unmarshal(b, &ipList); err != nil {
		panic(err)
	}

	return ipList.IpList
}

func GetAllNodesLocal(address string) ([]string, error) {
	addr := address + ":50052"
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	c := pb.NewOperationsClient(conn)

	res, err := c.GetAllNodes(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		return nil, err
	}

	return res.GetAddresses(), nil
}

func contactServer(ip string) (op.OperationsClient, *grpc.ClientConn, error) {
	addr := ip + ":50051"
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, nil, err
	}

	c := op.NewOperationsClient(conn)
	return c, conn, nil
}

func Ping(ip string) (time.Duration, error) {
	c, _, err := contactServer(ip)
	if err != nil {
		return -1, err
	}

	start := time.Now()
	_, err = c.Ping(context.Background(), &op.PingMessage{})
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed, nil
}

func GetClosestNode(nodes []string) (string, error) {
	min := time.Hour * 2
	minNode := ""

	for _, node := range nodes {
		rtt, err := Ping(node)
		if err != nil {
			continue
		}

		if rtt < min {
			min = rtt
			minNode = node
		}

	}

	if minNode == "" {
		return "", errors.New("no node available")
	}

	return minNode, nil
}

func Connect(ip string) (*Connection, *grpc.ClientConn, error) {
	c, conn, err := contactServer(ip)
	if err != nil {
		return nil, nil, err
	}

	return &Connection{c: c}, conn, nil
}

func (conn Connection) Get(k []byte) ([][]byte, error) {
	ret, err := conn.c.Get(context.Background(), &op.Key{Key: k})
	if err != nil {
		return nil, err
	}

	return ret.GetValue(), nil
}

func (conn Connection) Put(k []byte, val [][]byte) error {
	_, err := conn.c.Put(context.Background(), &op.KeyValue{Key: k, Value: val})
	if err != nil {
		return err
	}

	return nil
}

func (conn Connection) Append(k []byte, val [][]byte) error {
	_, err := conn.c.Append(context.Background(), &op.KeyValue{Key: k, Value: val})
	if err != nil {
		return err
	}

	return nil
}

func (conn Connection) Del(k []byte) error {
	_, err := conn.c.Del(context.Background(), &op.Key{Key: k})
	if err != nil {
		return err
	}

	return nil
}

/*
func main() {

	allNodes := GetAllNodesLocal("192.168.1.146")

	closest, _ := GetClosestNode(allNodes)

	conn, _ := Connect(closest)

	input := make([][]byte, 1)
	input[0] = []byte("defa")
	err := conn.Put([]byte("a"), input)
	if err != nil {
		log.Error(err)
	}

	val, err := conn.Get([]byte("a"))
	if err != nil {
		log.Error(err)
	}

	err = conn.Append([]byte("a"), input)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(val)

	err = conn.Del([]byte("a"))
	if err != nil {
		log.Error(err)
	}

	val, err = conn.Get([]byte("a"))
	if err != nil {
		log.Error(err)
	}

	fmt.Println(val)
}*/
