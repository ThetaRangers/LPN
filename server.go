package main

import (
	db "SDCC/database"
	"SDCC/ipfs"
	pb "SDCC/operations"
	"SDCC/utils"
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	port = ":50051"
	mask = "172.17.0.0/24"
	n    = 2 //Replication nodes
)

type Config struct {
	Port           int
	Seed           int64
	BootstrapPeers addrList
	TestMode       bool
}

type addrList []multiaddr.Multiaddr

var neighbours []string

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

type NullValidator struct{}

// Validate always returns success
func (nv NullValidator) Validate(key string, value []byte) error {
	log.Printf("NullValidator Validate: %s - %s", key, string(value))
	return nil
}

func ContactServer(ip string) (pb.OperationsClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewOperationsClient(conn)

	// Contact the server and print out its response.
	/*
		_, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
	*/
	return c, conn
}

// Select always selects the first record
func (nv NullValidator) Select(key string, values [][]byte) (int, error) {
	strs := make([]string, len(values))
	for i := 0; i < len(values); i++ {
		strs[i] = string(values[i])
	}
	log.Printf("NullValidator Select: %s - %v", key, strs)

	return 0, nil
}

var database db.Database
var ip net.IP
var address string

type server struct {
	pb.UnimplementedOperationsServer
}

func (s *server) Get(ctx context.Context, in *pb.Key) (*pb.Value, error) {
	//Request from the client
	log.Printf("Received: Get(%v)", in.GetKey())

	key := string(in.GetKey())
	value, err := kdht.GetValue(ctx, key)

	if err != nil {
		if err == routing.ErrNotFound {
			//Not found in the dht
			return &pb.Value{Value: [][]byte{}}, nil
		}
		return nil, err
	}

	remoteIp := string(value)

	if remoteIp != address {
		//Key present in this node
		//TODO handle failures
		c, _ := ContactServer(remoteIp)

		log.Println("Found key at  ", remoteIp, " connecting...")
		result, err := c.Get(ctx, &pb.Key{Key: in.GetKey(), Client: false})
		if err != nil {
			log.Fatal(err)
		}

		return result, nil
	} else {
		return &pb.Value{Value: database.Get(in.GetKey())}, nil
	}
}

func (s *server) Put(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	if in.GetClient() {
		log.Printf("Received: client Put(%v, %v)", in.GetKey(), in.GetValue())
		key := string(in.GetKey())
		//Check where is stored
		value, err := kdht.GetValue(ctx, key)
		if err != nil {
			if err == routing.ErrNotFound {
				log.Println("Not found responsible node, putting in local db....")
				// Not found in the dht
				database.Put(in.GetKey(), in.GetValue())

				//Set
				err := kdht.PutValue(ctx, string(in.GetKey()), []byte(address))
				if err != nil {
					return &pb.Ack{Msg: "Err"}, err
				}

				return &pb.Ack{Msg: "Ok"}, nil
			}

			log.Println("Error in put")
			return &pb.Ack{Msg: "Err"}, err
		}

		//Found in the dh
		remoteIp := string(value)

		if remoteIp != address {
			//Key present in this node
			log.Println("Found key at  ", remoteIp, " connecting...")
			c, _ := ContactServer(remoteIp)
			_, err := c.Put(ctx, &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue(), Client: false})
			if err != nil {
				return &pb.Ack{Msg: "Connection error"}, nil
			}

			return &pb.Ack{Msg: "Ok"}, nil
		}

	}

	log.Printf("PUTTING IN LOCAL DB: %s:%s\n", string(in.GetKey()), string(in.GetValue()))

	database.Put(in.GetKey(), in.GetValue())

	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Append(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	key := string(in.GetKey())

	if in.GetClient() {
		//Check where is stored
		value, err := kdht.GetValue(ctx, key)

		if err != nil {
			if err == routing.ErrNotFound {
				//Not found in the dht
				database.Put(in.GetKey(), in.GetValue())

				//Set
				err := kdht.PutValue(ctx, string(in.GetKey()), []byte(ip.String()))
				if err != nil {
					return nil, err
				}

				return &pb.Ack{Msg: "Ok"}, nil
			}

			return nil, err
		}

		remoteIp := string(value)

		//Found in the dht
		if remoteIp != address {
			//Connect to remote ip
			log.Println("Found key at  ", remoteIp, " connecting...")
			c, _ := ContactServer(remoteIp)
			_, err := c.Append(ctx, &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue(), Client: false})
			if err != nil {
				return &pb.Ack{Msg: "Connection error"}, nil
			}
		}
	}

	database.Append(in.GetKey(), in.GetValue())
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Del(ctx context.Context, in *pb.Key) (*pb.Ack, error) {
	key := string(in.GetKey())

	if in.GetClient() {
		//Delete in the DHT
		value, err := kdht.GetValue(ctx, key)
		if err != nil {
			if err == routing.ErrNotFound {
				// Not found in the dht
				//Can return
				return &pb.Ack{Msg: "Ok"}, nil
			}

			return &pb.Ack{Msg: "Err"}, err
		}

		remoteIp := string(value)

		//Found in the dht
		if remoteIp != address {
			log.Println("Found key at  ", remoteIp, " connecting...")
			c, _ := ContactServer(remoteIp)
			_, err := c.Del(ctx, &pb.Key{Key: []byte("abc"), Client: false})
			if err != nil {
				return &pb.Ack{Msg: "Connection error"}, nil
			}
		}
	}

	database.Del(in.GetKey())

	err := kdht.PutValue(ctx, key, []byte(""))
	if err != nil {
		return nil, err
	}
	//TODO do delete
	return &pb.Ack{Msg: "Ok"}, nil
}

func ContainsNetwork(mask string, ip net.IP) (bool, error) {
	_, subnet, err := net.ParseCIDR(mask)
	if err != nil {
		return false, err
	}
	return subnet.Contains(ip), err
}

func getIpFromPeerAddr(peers []multiaddr.Multiaddr) string {

	var parts []string
	for i := 0; i < len(peers); i++ {
		parts = strings.Split(peers[i].String(), "/")
		ipString := parts[2]
		//TODO CHANGE
		if !strings.Contains(ipString, "127.0.0") {
			return parts[2]
		}
		//ipres, _, _ := net.ParseCIDR(parts[2])
		/*res, _ := ContainsNetwork(mask, ipres)
		 if res {
			 return parts[2]
		 }*/
	}

	return parts[2]
}

func init() {
	database = utils.GetConfiguration().Database
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func sliceDiff(sl1, sl2 []string) []string {
	//sl1 items not contained in sl2
	var diff []string
	for _, ip := range sl1 {
		_, found := find(sl2, ip)
		if !found {
			diff = append(diff, ip)
		}
	}
	return diff
}

func neighboursPeers(ctx context.Context, nodeId string) {

	for true {
		//Find neighbours
		newNeighbours := make([]string, 0)
		peers, _ := kdht.GetClosestPeers(ctx, nodeId)
		fmt.Println(kdht.RoutingTable().ListPeers())
		if len(peers) > 0 {
			for i := 0; i < len(peers); i++ {
				val, _ := kdht.FindPeer(ctx, peers[i])

				peerIp := getIpFromPeerAddr(val.Addrs)
				newNeighbours = append(newNeighbours, peerIp)
				//newNeighbours[i] = peerIpr
			}
		}
		/*fmt.Println("Old", neighbours)

		sort.Strings(newNeighbours)
		fmt.Println("New", newNeighbours)*/

		//Find difference in the neighbours
		newPeers := sliceDiff(neighbours, newNeighbours)
		oldPeers := sliceDiff(newNeighbours, neighbours)

		fmt.Println("old neighbours ", neighbours)
		fmt.Println("new neighbours ", newNeighbours)

		fmt.Println("new peers ", newPeers)
		fmt.Println("old peers ", oldPeers)

		time.Sleep(5 * time.Second)
	}
}

var kdht *dht.IpfsDHT

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOperationsServer(s, &server{})

	neighbours = make([]string, n)

	bootstrap := os.Getenv("BOOTSTRAP_PEERS")
	if len(bootstrap) != 0 {
		log.Println("Found bootstrapp peer at ", bootstrap)
	}

	//Get ip address
	ifaces, err := net.Interfaces()
	// handle err
	for _, i := range ifaces {
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
				address = fmt.Sprintf("%s%s", ip, port)
				log.Printf("IP: %s", ip)
				break
			}
		}
	}

	// Joining the DHT
	config := Config{}
	flag.Int64Var(&config.Seed, "seed", 0, "Seed value for generating a PeerID, 0 is random")

	//For debugging
	if len(bootstrap) == 0 {
		flag.Var(&config.BootstrapPeers, "peer", "Peer multiaddress for peer discovery")
	} else {
		//addr, _ := multiaddr.NewMultiaddr(bootstrap)
		config.BootstrapPeers.Set(bootstrap)
	}

	flag.IntVar(&config.Port, "port", 0, "")
	flag.Parse()

	ctx := context.Background()

	h, err := ipfs.NewHost(ctx, config.Seed, config.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Host ID: %s", h.ID().Pretty())
	log.Printf("DHT addresses:")
	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID().Pretty())
	}

	kdht, err = ipfs.NewDHT(ctx, h, config.BootstrapPeers)
	kdht.Validator = NullValidator{}

	go neighboursPeers(ctx, h.ID().Pretty())

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
