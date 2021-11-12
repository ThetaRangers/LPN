package main

import (
	pb "SDCC/registerServer"
	"SDCC/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type RegisterStruct struct {
	Ip       string
	NodeId   string
	Attached bool
	Cluster  int
}

var list = make([]RegisterStruct, 0)
var maxCluster = -1
var lock sync.Mutex

type serverRegister struct {
	pb.UnimplementedOperationsServer
}

func getClusterNodes(id int, ip string) ([]string, []string) {
	tmpIp := make([]string, 0)
	tmpNodeId := make([]string, 0)

	for _, k := range list {
		if k.Cluster == id && k.Ip != ip {
			tmpIp = append(tmpIp, k.Ip)
			tmpNodeId = append(tmpNodeId, k.NodeId)
		}
	}

	return tmpIp, tmpNodeId
}

func getExistingCluster() ([]string, []string) {
	tmpIp := make([]string, 0)
	tmpNodeId := make([]string, 0)

	for _, k := range list {
		if k.Cluster == maxCluster {
			tmpIp = append(tmpIp, k.Ip)
			tmpNodeId = append(tmpNodeId, k.NodeId)
		}
	}

	return tmpIp, tmpNodeId
}

func isRegistered(ip string) (int, bool) {
	for index, v := range list {
		if v.Ip == ip {
			return index, true
		}
	}

	return -1, false
}

func clusterPossible() bool {
	count := 0
	for _, v := range list {
		if v.Cluster == -1 || v.Attached {
			count++
		}
	}

	if count == (utils.N - 1) {
		return true
	}

	return false
}

func createCluster() ([]string, []string, int) {
	tmpIp := make([]string, 0)
	tmpNodeId := make([]string, 0)

	maxCluster = maxCluster + 1

	for index, v := range list {
		if v.Cluster == -1 || v.Attached {
			tmpIp = append(tmpIp, v.Ip)
			tmpNodeId = append(tmpNodeId, v.NodeId)

			list[index].Cluster = maxCluster
		}
	}

	fmt.Println(tmpIp)
	return tmpIp, tmpNodeId, maxCluster
}

func (s *serverRegister) Register(ctx context.Context, in *pb.RegisterMessage) (*pb.Cluster, error) {
	lock.Lock()

	ip := in.GetIp()
	nodeId := in.GetNodeId()

	i, present := isRegistered(ip)
	if present {
		// If it is present in the network
		// Update NodeID
		list[i].NodeId = nodeId
		addresses, ids := getClusterNodes(list[i].Cluster, ip)

		lock.Unlock()
		return &pb.Cluster{Addresses: addresses, NodeIdS: ids, Crashed: true}, nil
	}

	if len(list) < (utils.N - 1) {
		// If not enough
		list = append(list, RegisterStruct{Ip: ip, NodeId: nodeId, Cluster: -1, Attached: false})

		lock.Unlock()
		return &pb.Cluster{Addresses: make([]string, 0), NodeIdS: make([]string, 0), Crashed: false}, nil
	} else if clusterPossible() {
		// Create new cluster
		addresses, ids, clusterId := createCluster()

		list = append(list, RegisterStruct{Ip: ip, NodeId: nodeId, Cluster: clusterId, Attached: false})

		lock.Unlock()
		return &pb.Cluster{Addresses: addresses, NodeIdS: ids, Crashed: false}, nil
	} else {
		// Attach
		list = append(list, RegisterStruct{Ip: ip, NodeId: nodeId, Cluster: maxCluster, Attached: true})
		addresses, ids := getExistingCluster()

		lock.Unlock()
		return &pb.Cluster{Addresses: addresses, NodeIdS: ids, Crashed: false}, nil
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOperationsServer(s, &serverRegister{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
