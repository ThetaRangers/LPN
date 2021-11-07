package replication

import (
	"SDCC/database"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"net"
	"os"
	"strings"
	"time"
)

const (
	cacheSize    = 512
	timeout      = 10 * time.Second
	maxPool      = 3
	raftPort     = 2500
	applyTimeout = 500 * time.Millisecond
)

type RaftStruct struct {
	RaftNode *raft.Raft
}

func (r RaftStruct) AddNode(ip string) error {
	address := fmt.Sprintf("%s:%d", ip, raftPort)

	f := r.RaftNode.AddVoter(raft.ServerID(ip), raft.ServerAddress(address), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}

	return nil
}

func (r RaftStruct) GetLeader() string {
	addr := r.RaftNode.Leader()

	s := strings.Split(string(addr), ":")

	return s[0]
}

func (r RaftStruct) Put(key []byte, value [][]byte) error {
	payload := CommandPayload{
		Operation: "PUT",
		Key:       key,
		Value:     value,
	}

	err := r.apply(payload)
	if err != nil {
		return err
	}

	return nil
}

func (r RaftStruct) Del(key []byte) error {
	payload := CommandPayload{
		Operation: "DELETE",
		Key:       key,
		Value:     nil,
	}

	err := r.apply(payload)
	if err != nil {
		return err
	}

	return nil
}

func (r RaftStruct) apply(payload CommandPayload) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	applyFuture := r.RaftNode.Apply(data, applyTimeout)
	if err := applyFuture.Error(); err != nil {
		return err
	}

	_, ok := applyFuture.Response().(*ApplyResponse)
	if !ok {
		return err
	}

	return nil
}

func (r RaftStruct) AddNodes(addresses []string) error {
	for _, addr := range addresses {
		err := r.AddNode(addr)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitializeRaft(ip string, db database.Database) RaftStruct {
	raftConf := raft.DefaultConfig()
	nodeId := ip
	raftConf.LocalID = raft.ServerID(nodeId)

	raftConf.SnapshotThreshold = 1024

	fsmStore := NewFSM(db)

	tcpTimeout := timeout
	var raftBinAddr = fmt.Sprintf("%s:%d", ip, raftPort)

	store, err := raftboltdb.NewBoltStore("raft-data/boltstore")
	if err != nil {
		panic(err)
	}

	cacheStore, err := raft.NewLogCache(cacheSize, store)
	if err != nil {
		panic(err)
	}

	snapshotStore, err := raft.NewFileSnapshotStore("raft-data/snapshot", 2, os.Stdout)
	if err != nil {
		panic(err)
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", raftBinAddr)
	if err != nil {
		panic(err)
	}

	transport, err := raft.NewTCPTransport(raftBinAddr, tcpAddr, maxPool, tcpTimeout, os.Stdout)
	if err != nil {
		panic(err)
	}

	raftNode, err := raft.NewRaft(raftConf, fsmStore, cacheStore, store, snapshotStore, transport)
	if err != nil {
		panic(err)
	}

	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      raftConf.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}

	raftNode.BootstrapCluster(configuration)

	r := RaftStruct{RaftNode: raftNode}

	return r
}
