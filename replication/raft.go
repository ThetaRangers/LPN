package replication

import (
	"SDCC/database"
	"SDCC/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"log"
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

func (r RaftStruct) Append(key []byte, value [][]byte) ([][]byte, bool, error) {
	payload := CommandPayload{
		Operation: "APPEND",
		Key:       key,
		Value:     value,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return [][]byte{}, false, err
	}

	applyFuture := r.RaftNode.Apply(data, applyTimeout)
	if err := applyFuture.Error(); err != nil {
		return [][]byte{}, false, err
	}

	response, ok := applyFuture.Response().(*AppendResponse)
	if !ok {
		return [][]byte{}, false, errors.New("bad response")
	} else {
		if response.Error != nil {
			return [][]byte{}, false, response.Error
		} else {
			if response.ToBeOffloaded {
				return response.Value, true, nil
			}
			return response.Value, false, nil
		}
	}

}

func (r RaftStruct) Join(ip string) error {
	payload := CommandPayload{
		Operation: "JOIN",
		Key:       []byte(ip),
	}

	err := r.apply(payload)
	if err != nil {
		return err
	}

	return nil
}

func (r RaftStruct) Leave(ip string) error {
	payload := CommandPayload{
		Operation: "LEAVE",
		Key:       []byte(ip),
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

func (r RaftStruct) AddNode(ip string) error {
	address := fmt.Sprintf("%s:%d", ip, raftPort)

	f := r.RaftNode.AddVoter(raft.ServerID(ip), raft.ServerAddress(address), 0, 0)
	if f.Error() != nil {
		return f.Error()
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

func (r RaftStruct) RemoveNode(ip string) error {
	f := r.RaftNode.RemoveServer(raft.ServerID(ip), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}

	return nil
}

func ReInitializeRaft(ip string, db *database.Database, cluster *utils.ClusterRoutine, dht *DhtRoutine) *RaftStruct {
	err := os.RemoveAll("data/raft-data")
	if err != nil {
		log.Fatal(err)
	}

	return InitializeRaft(ip, db, cluster, dht)
}

func InitializeRaft(ip string, db *database.Database, cluster *utils.ClusterRoutine, dht *DhtRoutine) *RaftStruct {
	err := os.RemoveAll("data/raft-data")
	if err != nil {
		panic(err)
	}

	err = os.Mkdir("data/raft-data", 0755)
	if err != nil {
		log.Println(err)
	}

	raftConf := raft.DefaultConfig()
	nodeId := ip
	raftConf.LocalID = raft.ServerID(nodeId)

	raftConf.SnapshotThreshold = 1024

	fsmStore := NewFSM(db, cluster, dht)

	tcpTimeout := timeout
	var raftBinAddr = fmt.Sprintf("%s:%d", ip, raftPort)

	ldb, err := raftboltdb.NewBoltStore("data/raft-data/logstore")
	if err != nil {
		log.Println("boltdb.NewBoltStore:", err)
	}

	sdb, err := raftboltdb.NewBoltStore("data/raft-data/stablestore")
	if err != nil {
		log.Println("boltdb.NewBoltStore:", err)
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

	raftNode, err := raft.NewRaft(raftConf, fsmStore, ldb, sdb, snapshotStore, transport)
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

	return &RaftStruct{RaftNode: raftNode}
}
