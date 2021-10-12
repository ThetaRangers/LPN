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
)

const (
	port = ":50051"
	mask = "127.0.0.1/8"
)

type Config struct {
	Port           int
	Seed           int64
	BootstrapPeers addrList
	TestMode       bool
}

type addrList []multiaddr.Multiaddr

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

type server struct {
	pb.UnimplementedOperationsServer
}

func (s *server) Get(ctx context.Context, in *pb.Key) (*pb.Value, error) {
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
	fmt.Println("Found at: ", remoteIp)

	if remoteIp == string(ip) {
		//Key present in this node
		return &pb.Value{Value: database.Get(in.GetKey())}, nil
	} else {
		//TODO handle connections
		//TODO handle failures
		return &pb.Value{Value: database.Get(in.GetKey())}, nil
	}
}

func (s *server) Put(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: Put(%v, %v)", in.GetKey(), in.GetValue())

	key := string(in.GetKey())
	//Check where is stored
	value, err := kdht.GetValue(ctx, key)
	if err != nil {
		if err == routing.ErrNotFound {
			// Not found in the dht
			database.Put(in.GetKey(), in.GetValue())

			//Set
			err := kdht.PutValue(ctx, string(in.GetKey()), []byte(ip.String()))
			if err != nil {
				return &pb.Ack{Msg: "Err"}, err
			}

			return &pb.Ack{Msg: "Ok"}, nil
		}

		return &pb.Ack{Msg: "Err"}, err
	}

	fmt.Println("Got value ", value)

	//Found in the dh
	remoteip := string(value)
	fmt.Println("Found at: ", remoteip)

	if remoteip == string(ip) {
		//Key present in this node
		database.Put(in.GetKey(), in.GetValue())
	} else {
		//TODO contact remote server
	}

	log.Printf("Found")

	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Append(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: Append(%v, %v)", in.GetKey(), in.GetValue())
	key := string(in.GetKey())
	//Check where is stored
	value, err := kdht.GetValue(ctx, key)
	if err != nil {
		if err == routing.ErrNotFound {
			//Not found in the dht
			database.Put(in.GetKey(), in.GetValue())

			//Set
			//TODO
			err := kdht.PutValue(ctx, string(in.GetKey()), []byte(ip.String()))
			if err != nil {
				return nil, err
			}

			return &pb.Ack{Msg: "Ok"}, nil
		}

		return nil, err
	}

	remoteip := string(value)

	//Found in the dht
	if remoteip == string(ip) {
		//Key present in this node
		database.Append(in.GetKey(), in.GetValue())
	} else {
		//TODO contact remote server
	}

	log.Printf("Found")

	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Del(ctx context.Context, in *pb.Key) (*pb.Ack, error) {
	log.Printf("Received: Del(%v)", in.GetKey())

	key := string(in.GetKey())
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

	remoteip := string(value)

	//Found in the dht
	if remoteip == string(ip) {
		//Key present in this node
		database.Del(in.GetKey())
	} else {
		//TODO contact remote server

	}

	database.Del(in.GetKey())
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

func init() {
	database = utils.GetConfiguration().Database
}

var kdht *dht.IpfsDHT

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOperationsServer(s, &server{})

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
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
