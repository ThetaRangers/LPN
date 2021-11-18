package client

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

type Connection struct {
	c op.OperationsClient
}

type IpListStruct struct {
	IpList []string `json:"ipList"`
}

func GetAllNodes(url string) []string {
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
	min := time.Hour
	minNode := ""

	channels := make([]chan time.Duration, 0)

	for i, k := range nodes {
		channels = append(channels, make(chan time.Duration))

		k := k
		i := i
		go func() {
			rtt, err := Ping(k)
			if err != nil {
				channels[i] <- -1
			} else {
				channels[i] <- rtt
			}
		}()
	}

	for i, c := range channels {
		rtt := <-c
		if rtt < min && rtt > 0 {
			min = rtt
			minNode = nodes[i]
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
