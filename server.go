package main

import (
	"SDCC/cloud"
	db "SDCC/database"
	"SDCC/dht"
	"SDCC/metadata"
	"SDCC/migration"
	pb "SDCC/operations"
	"SDCC/replication"
	"SDCC/utils"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	port       = ":50051"
	mask       = "172.17.0.0/24"
	regService = 1
	threshold  = 300
)

var retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "operations.Operations"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`

type Config struct {
	Port           int
	Seed           int64
	BootstrapPeers addrList
	TestMode       bool
}

type server struct {
	pb.UnimplementedOperationsServer
}

type addrList []multiaddr.Multiaddr

var replicaSet []string
var cluster utils.ClusterRoutine

var database db.Database
var ip net.IP
var address string
var channel chan migration.KeyOp
var raftN *replication.RaftStruct
var dhtRoutine replication.DhtRoutine
var kdht *dht.KDht
var config = Config{}
var h host.Host

var keyDb = metadata.GetKeyDb()

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}

func ContactServer(ip string) (pb.OperationsClient, *grpc.ClientConn, error) {
	addr := ip + ":50051"
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		return nil, nil, err
	}

	c := pb.NewOperationsClient(conn)
	return c, conn, nil
}

func getAliveReplica(ctx context.Context, ip string) (pb.OperationsClient, *grpc.ClientConn, error) {
	c, conn, err := ContactServer(ip)
	if err != nil {
		replicas, err2 := kdht.GetCluster(ctx, ip)
		if err2 != nil {
			return nil, nil, err2
		}
		// ip is part of replicas, iterating this way may cause retrying to connect to it
		for _, replica := range replicas {
			c, conn, err = ContactServer(replica)
			if err == nil {
				return c, conn, nil
			}
		}
		return nil, nil, errors.New("no alive replica found")
	} else {
		return c, conn, nil
	}
}

// Get rpc function called to retrieve a value, if the value is not found in the local DB, the responsible node is
// searched on the DHT and queried. If no node can be found an empty value is returned with no error. If an error
// occurred an empty value is returned with the error. If the value is correctly found, it is returned with no error
func (s *server) Get(ctx context.Context, in *pb.Key) (*pb.Value, error) {
	//Request from the client
	log.Printf("Received: Get(%v)", in.GetKey())
	key := string(in.GetKey())

	value, onTheCloud, err := kdht.GetValue(ctx, key)
	if onTheCloud {
		//TODO GET FROM CLOUD
	} else if err == routing.ErrNotFound || err == nil && len(value) == 0 {
		//Not found in the dht or deleted
		return &pb.Value{Value: [][]byte{}}, nil
	} else if err != nil {
		return &pb.Value{Value: [][]byte{}}, err
	}

	if !cluster.Contains(value) {
		channel <- migration.KeyOp{Key: key, Op: migration.ReadOperation, Mode: migration.External}

		// Try node list
		c, _, err := getAliveReplica(ctx, value)
		if err != nil {
			return &pb.Value{Value: [][]byte{}}, err
		}

		result, err := c.GetInternal(ctx, &pb.Key{Key: in.GetKey()})

		return result, nil
		//return &pb.Value{Value: [][]byte{}}, errors.New("All replicas down")

	} else {
		value, _ := database.Get(in.GetKey())
		channel <- migration.KeyOp{Key: key, Op: migration.ReadOperation, Mode: migration.Master}

		return &pb.Value{Value: value}, nil
	}
}

// GetInternal internal function called by other nodes to retrieve an information
func (s *server) GetInternal(ctx context.Context, in *pb.Key) (*pb.Value, error) {
	val, _ := database.Get(in.GetKey())
	return &pb.Value{Value: val}, nil
}

// Put rpc function called to store a value on the responsible node. If no responsible node is found, the current node
// becomes the responsible.
func (s *server) Put(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: client Put(%v, %v)", in.GetKey(), in.GetValue())

	ctxDht := context.Background()
	key := string(in.GetKey())

	var offload bool
	//Check where is stored
	value, _, err := kdht.GetValue(ctxDht, key)
	if err == routing.ErrNotFound || err == nil && len(value) == 0 {
		if utils.GetSize(in.GetValue()) > threshold {
			// TODO OFFLOAD TO THE CLOUD
			address = "onTheCloud"
			offload = true
		} else {
			log.Println("Not found responsible node, putting in local db....")
			channel <- migration.KeyOp{Key: key, Op: migration.WriteOperation, Mode: migration.Master}
			// Not found in the dht or deleted

			keyDb.PutKey(key)
			leader := raftN.GetLeader()
			if leader == ip.String() {
				err = raftN.Put(in.GetKey(), in.GetValue())
			} else {
				c, _, _ := ContactServer(leader)

				_, err = c.PutInternal(context.Background(), &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
			}

			if err != nil {
				return &pb.Ack{Msg: "Err"}, err
			}

			offload = false
		}

		err := kdht.PutValue(ctxDht, string(in.GetKey()), address, offload)
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		return &pb.Ack{Msg: "Ok"}, nil
	} else if err != nil {
		return &pb.Ack{Msg: "Err"}, err
	}

	//Found in the dht
	// If not in Raft cluster
	if !cluster.Contains(value) {
		//Connect to remote ip
		channel <- migration.KeyOp{Key: key, Op: migration.WriteOperation, Mode: migration.External}
		c, _, err := getAliveReplica(ctx, value)
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		_, err = c.PutInternal(ctx, &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		return &pb.Ack{Msg: "Ok"}, nil
	} else {
		channel <- migration.KeyOp{Key: key, Op: migration.WriteOperation, Mode: migration.Master}

		leader := raftN.GetLeader()
		if leader == ip.String() {
			err = raftN.Put(in.GetKey(), in.GetValue())
		} else {
			c, _, _ := ContactServer(leader)

			_, err = c.PutInternal(context.Background(), &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
		}

		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}
	}

	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) PutInternal(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	var err error

	leader := raftN.GetLeader()
	if leader == ip.String() {
		err = raftN.Put(in.GetKey(), in.GetValue())
	} else {
		c, _, _ := ContactServer(leader)

		_, err = c.PutInternal(context.Background(), &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
	}
	if err != nil {
		return &pb.Ack{Msg: "Err"}, err
	}

	return &pb.Ack{Msg: "Ok"}, nil
}

// Append i i no green pass
func (s *server) Append(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	key := string(in.GetKey())

	//Check where is stored
	var offload bool
	//Check where is stored
	value, _, err := kdht.GetValue(ctx, key)
	if err == routing.ErrNotFound || err == nil && len(value) == 0 {
		if utils.GetSize(in.GetValue()) > threshold {
			// TODO OFFLOAD TO THE CLOUD
			address = "onTheCloud"
			offload = true
		} else {
			log.Println("Not found responsible node, putting in local db....")
			channel <- migration.KeyOp{Key: key, Op: migration.WriteOperation, Mode: migration.Master}
			// Not found in the dht or deleted

			keyDb.PutKey(key)
			leader := raftN.GetLeader()
			if leader == ip.String() {
				err = raftN.Put(in.GetKey(), in.GetValue())
			} else {
				c, _, _ := ContactServer(leader)

				_, err = c.PutInternal(context.Background(), &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
			}

			if err != nil {
				return &pb.Ack{Msg: "Err"}, err
			}

			offload = false
		}

		err := kdht.PutValue(ctx, string(in.GetKey()), address, offload)
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		return &pb.Ack{Msg: "Ok"}, nil
	} else if err != nil {
		return nil, err
	}

	//Found in the dht
	if !cluster.Contains(value) {
		//Connect to remote ip
		channel <- migration.KeyOp{Key: key, Op: migration.WriteOperation, Mode: migration.External}

		c, _, err := getAliveReplica(ctx, value)
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		_, err = c.AppendInternal(ctx, &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
		if err != nil {
			return &pb.Ack{Msg: "Connection error"}, nil
		}
	} else {
		channel <- migration.KeyOp{Key: key, Op: migration.WriteOperation, Mode: migration.Master}

		leader := raftN.GetLeader()
		if leader == ip.String() {
			err = raftN.Append(in.GetKey(), in.GetValue())
		} else {
			c, _, _ := ContactServer(leader)

			_, err = c.AppendInternal(context.Background(), &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
		}
	}

	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) AppendInternal(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	var err error

	leader := raftN.GetLeader()
	if leader == ip.String() {
		err = raftN.Append(in.GetKey(), in.GetValue())
	} else {
		c, _, _ := ContactServer(leader)

		_, err = c.AppendInternal(context.Background(), &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
	}
	if err != nil {
		return &pb.Ack{Msg: "Err"}, err
	}

	return &pb.Ack{Msg: "Ok"}, nil
}

// Del function to delete
func (s *server) Del(ctx context.Context, in *pb.Key) (*pb.Ack, error) {
	key := string(in.GetKey())
	// Get responsible from dht
	value, onTheCloud, err := kdht.GetValue(ctx, key)
	if onTheCloud {
		// TODO DELETE KEY FROM CLOUD
	} else if err == routing.ErrNotFound || err == nil && len(value) == 0 {
		// Not found in the dht or already deleted
		return &pb.Ack{Msg: "Ok"}, nil
	} else if err != nil {
		return &pb.Ack{Msg: "Err"}, err
	}

	if !cluster.Contains(value) {
		c, _, err := getAliveReplica(ctx, value)
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		return c.DelInternal(ctx, in)
	} else {
		return s.DelInternal(ctx, in)
	}
}

func (s *server) DelInternal(ctx context.Context, in *pb.Key) (*pb.Ack, error) {
	var err error
	leader := raftN.GetLeader()
	if leader == address {
		err = raftN.Del(in.GetKey())
	} else {
		c, _, _ := ContactServer(leader)
		_, err = c.DelInternal(context.Background(), &pb.Key{Key: in.GetKey()})
	}
	if err != nil {
		return &pb.Ack{Msg: "Err"}, err
	}

	err = kdht.PutValue(ctx, string(in.GetKey()), "", false)
	if err != nil {
		return &pb.Ack{Msg: "Err"}, err
	}
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Migration(ctx context.Context, in *pb.KeyCost) (*pb.Outcome, error) {
	keyBytes := in.GetKey()
	k := string(keyBytes)

	log.Println("Received migration request for ", k)
	cost := uint64(migration.GetCostMaster(k, time.Now()))

	if cost < in.Cost {
		// Do migration
		value, err := database.Get(keyBytes)
		if err != nil {
			return &pb.Outcome{Out: false}, nil
		}

		err = raftN.Del(keyBytes)
		if err != nil {
			return &pb.Outcome{Out: false}, nil
		}

		migration.SetExported(k)

		return &pb.Outcome{Out: true, Value: value}, nil
	} else {
		// Do nothing
		log.Println("Migration refused")
		return &pb.Outcome{Out: false}, nil
	}
}

// Join Raft leader calls this to ask to Join the cluster
func (s *server) Join(ctx context.Context, in *pb.JoinMessage) (*pb.Ack, error) {
	if cluster.Len() != 0 {
		// If node is part of a cluster and needs to transfer
		// TODO transfer cluster

		log.Println("Leaving old cluster: ", cluster)
		leader := raftN.GetLeader()
		log.Println("Leader is ", leader)
		c, _, err := ContactServer(leader)
		if err != nil {
			log.Println("ERROR", err)
			return nil, err
		}

		_, err = c.LeaveCluster(context.Background(), &pb.RequestJoinMessage{Ip: ip.String()})
		if err != nil {
			log.Println("ERROR", err)
			return nil, err
		}

		f := raftN.RaftNode.Shutdown()
		err = f.Error()
		if err != nil {
			log.Fatal("ERROR IN SHUTDOWN", err)
		}

		raftN = replication.ReInitializeRaft(ip.String(), &database, &cluster, &dhtRoutine)
	} else {
		var err error
		bootstrapNodes := in.GetBootstrap()
		for _, b := range bootstrapNodes {
			err := config.BootstrapPeers.Set(b)
			if err != nil {
				log.Println(err)
			}

		}

		kdht, err = dht.NewKDht(ctx, h, config.BootstrapPeers)
		if err != nil {
			log.Fatal(err)
		}
		dhtRoutine.SetDht(kdht)
	}

	return &pb.Ack{Msg: "Ok"}, nil
}

// LeaveCluster used to leave a cluster
func (s *server) LeaveCluster(ctx context.Context, in *pb.RequestJoinMessage) (*pb.Ack, error) {
	address := in.GetIp()

	err := raftN.Leave(address)
	if err != nil {
		return nil, err
	}

	err = raftN.RemoveNode(address)
	if err != nil {
		return nil, err
	}

	return &pb.Ack{Msg: "Ok"}, nil
}

// RequestJoin Function called to request a join in a cluster
func (s *server) RequestJoin(ctx context.Context, in *pb.RequestJoinMessage) (*pb.JoinMessage, error) {
	address := in.GetIp()

	// If this node is the master
	leader := raftN.GetLeader()
	if raftN.GetLeader() == ip.String() {
		err := raftN.Join(address)
		if err != nil {
			return nil, err
		}

		err = raftN.AddNode(address)
		if err != nil {
			return nil, err
		}
	} else {
		c, _, err := ContactServer(leader)
		if err != nil {
			return nil, err
		}

		_, err = c.RequestJoin(ctx, &pb.RequestJoinMessage{Ip: address})
		if err != nil {
			return nil, err
		}
	}

	return &pb.JoinMessage{}, nil
}

func (s *server) Ping(_ context.Context, _ *pb.PingMessage) (*pb.Ack, error) {
	return &pb.Ack{Msg: "Pong"}, nil
}

func ContainsNetwork(mask string, ip net.IP) (bool, error) {
	_, subnet, err := net.ParseCIDR(mask)
	if err != nil {
		return false, err
	}
	return subnet.Contains(ip), err
}

func init() {
	database = utils.GetConfiguration().Database
}

func houseKeeper(ctx context.Context) {
	var val [][]byte
	for {
		// Evaluate migration
		migrationKeys := migration.EvaluateMigration()

		for _, k := range migrationKeys {
			value, fromCloud, err := kdht.GetValue(ctx, k)
			if err != nil {
				continue
			}

			if fromCloud {
				continue
			} else {
				// Try to contact server
				c, _, err := ContactServer(value)
				if err != nil {
					continue
				}

				outcome, err := c.Migration(ctx, &pb.KeyCost{Key: []byte(k), Cost: uint64(migration.GetCostExternal(k, time.Now()))})
				if err != nil || !outcome.Out {
					continue
				}
				val = outcome.Value
			}

			// Do migration
			migration.SetMigrated(k)
			err = raftN.Put([]byte(k), val)
			if err != nil {
				return
			}

			// Modify dht
			err = kdht.PutValue(ctx, k, ip.String(), false)
			if err != nil {
				// TODO ???
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func initializeHost(ctx context.Context) string {
	var addrString string
	var err error

	flag.Int64Var(&config.Seed, "seed", 0, "Seed value for generating a PeerID, 0 is random")

	h, err = dht.NewHost(ctx, config.Seed, config.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Host ID: %s", h.ID().Pretty())
	log.Printf("DHT addresses:")
	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID().Pretty())
		if strings.Contains(addr.String(), ip.String()) {
			addrString = addr.String()
		}
	}

	ipStr := fmt.Sprintf("%s/p2p/%s", addrString, h.ID().Pretty())

	return ipStr
}

func main() {
	//Get ip address
	iFaces, err := net.Interfaces()
	// handle err
	for _, i := range iFaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			check, err := ContainsNetwork(mask, ip)
			if err != nil {
				log.Panic(err)
			}

			if check {
				//address = fmt.Sprintf("%s%s", ip, port)
				address = ip.String()
				log.Printf("IP: %s", ip)
				break
			}
		}
	}

	// Initialize logging channel
	channel = make(chan migration.KeyOp, 200)
	go migration.ManagementThread(channel, utils.CostRead, utils.CostWrite, utils.MigrationWindowMinutes)

	// Initialize migration thread
	go houseKeeper(context.Background())

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOperationsServer(s, &server{})

	bootstrap := os.Getenv("BOOTSTRAP_PEERS")
	if len(bootstrap) != 0 {
		log.Println("Found bootstrap peer at ", bootstrap)
	}

	// Joining the DHT

	//For debugging
	if len(bootstrap) == 0 {
		//flag.Var(&config.BootstrapPeers, "peer", "Peer multiaddress for peer discovery")
	} else {
		//addr, _ := multiaddr.NewMultiaddr(bootstrap)
		//config.BootstrapPeers.Set(bootstrap)
	}

	//flag.IntVar(&config.Port, "port", 0, "")
	//flag.Parse()

	// Initialize Raft replication
	cluster = utils.NewClusterRoutine()
	dhtRoutine = replication.NewDhtRoutine(ip.String(), &cluster)
	raftN = replication.InitializeRaft(ip.String(), &database, &cluster, &dhtRoutine)
	time.Sleep(3 * time.Second)

	ctx := context.Background()
	ipStr := initializeHost(ctx)

	var registerCluster cloud.ReplicaSet
	if regService == 0 {
		// TODO maybe add support for ip6
		registerCluster = cloud.RegisterToTheNetwork(ip.String(), ipStr, utils.N, utils.AwsRegion)
	} else if regService == 1 {
		registerCluster = cloud.RegisterStub(ip.String(), ipStr, utils.N, utils.AwsRegion)
	}

	// Ready to start dht
	if len(registerCluster.IpList) > 0 {
		for j := 0; j < len(registerCluster.IpList); j++ {
			r := registerCluster.IpList[j]
			if r.IpString != ipStr {
				replicaSet = append(replicaSet, r.Ip)
				err := config.BootstrapPeers.Set(r.IpString)
				if err != nil {
					panic(err)
				}
			}
		}

		kdht, err = dht.NewKDht(ctx, h, config.BootstrapPeers)
		if err != nil {
			log.Fatal(err)
		}
		dhtRoutine.SetDht(kdht)
	}

	ipList := registerCluster.IpList

	if registerCluster.Crashed == 1 {
		// Node is crashed need to rejoin the cluster
		// TODO retry
		c, _, _ := ContactServer(ipList[0].Ip)
		_, err = c.RequestJoin(context.Background(), &pb.RequestJoinMessage{Ip: ip.String()})
		if err != nil {
			log.Fatal("ERROR IN JOINING", err)
		}

	} else if len(ipList) == (utils.N - 1) {
		// If master to a cluster
		log.Println("Replicas found: ", replicaSet)
		log.Println("Cluster: ", cluster)

		var clusterBootstrap []string
		var clusterAddresses []string
		for _, node := range ipList {
			clusterBootstrap = append(clusterBootstrap, node.IpString)
			clusterAddresses = append(clusterAddresses, node.Ip)
		}

		for _, node := range ipList {
			var err error
			var c pb.OperationsClient

			log.Printf("Asking %s to join", node.Ip)

			for err != nil || c == nil {
				c, _, err = ContactServer(node.Ip)
				log.Println("Failed to contact, trying again...")
			}

			_, err = c.Join(context.Background(), &pb.JoinMessage{Bootstrap: clusterBootstrap, Cluster: clusterAddresses})
			if err != nil {
				log.Fatal("ERROR in requesting join", err)
			}

			log.Printf("Adding %s to raft, leader is %s", node.Ip, raftN.GetLeader())

			err = raftN.AddNode(node.Ip)
			if err != nil {
				log.Fatal("ADDING NODE", err)
			}
			log.Printf("Done adding %s to raft", node.Ip)

		}
		raftN.Join(ip.String())
		for _, addr := range replicaSet {
			raftN.Join(addr)
		}

	} else if len(replicaSet) >= utils.N {
		// If external to the cluster
		target := ipList[0].Ip

		// TODO try to contact others
		log.Println("Requesting join on cluster ", replicaSet)
		c, _, _ := ContactServer(target)

		_, err := c.RequestJoin(context.Background(), &pb.RequestJoinMessage{Ip: ip.String()})
		if err != nil {
			log.Println("ERROR", err)
		}

	}

	go func() {
		/*var dab *badger.DB
		dab = database.(db.BadgerDB).Db*/
		for {
			/*log.Printf("Stats for : %s - %s", ip.String(), raftN.RaftNode.Stats())
			log.Println("____________________Begin printing db____________________")
			err := dab.View(func(txn *badger.Txn) error {
				opts := badger.DefaultIteratorOptions
				opts.PrefetchSize = 10
				it := txn.NewIterator(opts)
				defer it.Close()
				for it.Rewind(); it.Valid(); it.Next() {
					item := it.Item()
					k := item.Key()
					err := item.Value(func(v []byte) error {
						fmt.Printf("key=%s, value=%s\n", k, v)
						return nil
					})
					if err != nil {
						return err
					}
				}
				return nil
			})
			log.Println("____________________End printing db____________________")
			if err != nil {
				log.Println("There is a problem")
			}
			*/
			//log.Printf("Stats for : %s - %s", ip.String(), raftN.RaftNode.Stats()["latest_configuration"])

			log.Printf("Keys of %s: %s", ip.String(), keyDb.GetKeys())
			time.Sleep(15 * time.Second)
		}
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
