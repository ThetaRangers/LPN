package main

import (
	"SDCC/cloud"
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
)

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

var database db.Database
var ip net.IP
var address string

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
func (nv NullValidator) Validate(string, []byte) error {
	//log.Printf("NullValidator Validate: %s - %s", key, string(value))
	return nil
}

// Select always selects the first record
func (nv NullValidator) Select(string, [][]byte) (int, error) {
	/*
		strs := make([]string, len(values))
		for i := 0; i < len(values); i++ {
			strs[i] = string(values[i])
		}
		log.Printf("NullValidator Select: %s - %v", key, strs)
	*/

	return 0, nil
}

func ContactServer(ip string) (pb.OperationsClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	c := pb.NewOperationsClient(conn)
	return c, conn, nil
}

/*
 * La get anche se riguarda una replica, chiede il valore al master, così avviene una riconciliazione del valore
 */

// Get rpc function called to retrieve a value, if the value is not found in the local DB, the responsible node is
// searched on the DHT and queried. If no node can be found an empty value is returned with no error. If an error
// occurred an empty value is returned with the error. If the value is correctly found, it is returned with no error
/* func (s *server) Get(ctx context.Context, in *pb.Key) (*pb.Value, error) {
	//Request from the client
	log.Printf("Received: Get(%v)", in.GetKey())
	key := string(in.GetKey())
	value, err := kdht.GetValue(ctx, key)
	if err != nil {
		if err == routing.ErrNotFound {
			//Not found in the dht
			return &pb.Value{Value: [][]byte{}}, nil
		} else {
			return &pb.Value{Value: [][]byte{}}, err
		}
	}

	remoteIp := string(value) // TODO list values
	// bool replica = list.contains(me)

	if remoteIp != address {
		// Try node list
		//i := 0
		c, _, err := ContactServer(remoteIp)
		if err != nil {
			for {
				if replica {
					nuove elezioni
				} else {
					i++
					remoteIp = remoteIp[i]
					c, _, err := ContactServer(remoteIp)
			}
		}


		for {
			//if i > list.size() break;
			//TODO skip to next one in the list
			c, _, err := ContactServer(remoteIp)
			log.Println("Get ContactServer failure", err)
			if err != nil {
				// i++
				continue
			}

			result, err := c.GetInternal(ctx, &pb.Key{Key: in.GetKey()})
			if err != nil {
				// i++
				continue
			}
			return result, nil
		}
		//return &pb.Value{Value: [][]byte{}}, errors.New("All replicas down")

	} else {
		return &pb.Value{Value: database.Get(in.GetKey())}, nil
	}
} */

func (s *server) Get(ctx context.Context, in *pb.Key) (*pb.Value, error) {
	//Request from the client
	log.Printf("Received: Get(%v)", in.GetKey())
	key := string(in.GetKey())
	value, err := kdht.GetValue(ctx, key)
	if err != nil {
		if err == routing.ErrNotFound {
			//Not found in the dht
			return &pb.Value{Value: [][]byte{}}, nil
		} else {
			return &pb.Value{Value: [][]byte{}}, err
		}
	}

	remoteIp := string(value) // TODO list values

	if remoteIp != address {
		// Try node list
		//i := 0
		for {
			//if i > list.size() break;
			//TODO skip to next one in the list
			c, _, err := ContactServer(remoteIp)
			log.Println("Get ContactServer failure", err)
			if err != nil {
				// i++
				continue
			}

			result, err := c.GetInternal(ctx, &pb.Key{Key: in.GetKey()})
			if err != nil {
				// i++
				continue
			}
			return result, nil
		}
		//return &pb.Value{Value: [][]byte{}}, errors.New("All replicas down")

	} else {
		value, versionNum := database.Get(in.GetKey())
		fmt.Println("Version:", versionNum)
		return &pb.Value{Value: value}, nil
	}
}

// Put rpc function called to store a value on the responsible node. If no responsible node is found, the current node
// becomes the responsible.
func (s *server) Put(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	log.Printf("Received: client Put(%v, %v)", in.GetKey(), in.GetValue())

	ctxDht := context.Background()
	var dbInput [][]byte

	key := string(in.GetKey())
	//Check where is stored
	value, err := kdht.GetValue(ctxDht, key)
	if err != nil {
		if err == routing.ErrNotFound {
			log.Println("Not found responsible node, putting in local db....")
			// Not found in the dht
			dbInput = append(dbInput, in.GetValue())
			database.Put(in.GetKey(), dbInput)
			channel := make(chan bool)
			for i := 0; i < utils.Replicas; i++ {
				// Replicate as goroutine
				go func() {
					channel <- true
				}()
			}
			go func() {
				time.Sleep(utils.Timeout)
				channel <- false
			}()
			for i := 0; i < utils.WriteQuorum; i++ {
				done := <-channel
				if !done {
					// Timeout
				}
			}
			//Set
			err := kdht.PutValue(ctxDht, string(in.GetKey()), []byte(address))
			if err != nil {
				return &pb.Ack{Msg: "Err"}, err
			}

			return &pb.Ack{Msg: "Ok"}, nil
		} else {
			return &pb.Ack{Msg: "Err"}, err
		}
	}

	//Found in the dht
	remoteIp := string(value)

	if remoteIp != address {
		//Connect to remote ip
		//log.Println("Found key at  ", remoteIp, " connecting...")
		c, _, err := ContactServer(remoteIp)
		for err != nil {
			//TODO skip to next one in the list
			//c, _, err = ContactServer(remoteIp)
			log.Println("Put ContactServer failure", err)
			break
		}

		_, err = c.Put(ctx, &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
		if err != nil {
			return &pb.Ack{Msg: "Err"}, err
		}

		return &pb.Ack{Msg: "Ok"}, nil
	}

	dbInput = append(dbInput, in.GetValue())
	database.Put(in.GetKey(), dbInput)

	return &pb.Ack{Msg: "Ok"}, nil
}

// Append i i no green pass
func (s *server) Append(ctx context.Context, in *pb.KeyValue) (*pb.Ack, error) {
	key := string(in.GetKey())

	//Check where is stored
	value, err := kdht.GetValue(ctx, key)
	var dbInput [][]byte

	if err != nil {
		if err == routing.ErrNotFound {
			//Not found in the dht

			dbInput = append(dbInput, in.GetValue())
			database.Put(in.GetKey(), dbInput)

			//Set
			err := kdht.PutValue(ctx, string(in.GetKey()), []byte(ip.String()))
			if err != nil {
				return nil, err
			}

			return &pb.Ack{Msg: "Ok"}, nil
		} else {
			<-kdht.ForceRefresh()
		}

		return nil, err
	}

	remoteIp := string(value)

	//Found in the dht
	if remoteIp != address {
		//Connect to remote ip
		log.Println("Found key at  ", remoteIp, " connecting...")
		c, _, _ := ContactServer(remoteIp)
		_, err := c.Append(ctx, &pb.KeyValue{Key: in.GetKey(), Value: in.GetValue()})
		if err != nil {
			return &pb.Ack{Msg: "Connection error"}, nil
		}
	}

	database.Append(in.GetKey(), in.GetValue())
	return &pb.Ack{Msg: "Ok"}, nil
}

// Del function to delete
func (s *server) Del(ctx context.Context, in *pb.Key) (*pb.Ack, error) {
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
	} else {
		<-kdht.ForceRefresh()
	}

	remoteIp := string(value)

	//Found in the dht
	if remoteIp != address {
		log.Println("Found key at  ", remoteIp, " connecting...")
		c, _, _ := ContactServer(remoteIp)
		_, err := c.Del(ctx, &pb.Key{Key: []byte("abc")})
		if err != nil {
			return &pb.Ack{Msg: "Connection error"}, nil
		}
	}

	database.Del(in.GetKey())

	err = kdht.PutValue(ctx, key, []byte(""))
	if err != nil {
		return nil, err
	}
	//TODO do delete
	return &pb.Ack{Msg: "Ok"}, nil
}

func (s *server) Replicate(ctx context.Context, in *pb.KeyValueVersion) (*pb.Ack, error) {
	database.Put(in.GetKey(), in.GetValue(), in.GetVersion())
	return &pb.Ack{Msg: "Ok"}, nil
}

func callReplicate(ctx context.Context, ip string, key []byte, value [][]byte, version uint64) {
	c, _, _ := ContactServer(ip)
	ack, err := c.Replicate(ctx, &pb.KeyValueVersion{Key: key, Value: value, Version: version})
	if err != nil {
		return
	}
	if ack.GetMsg() != "Ok" {
		// TODO
	}
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

	replicaSet = cloud.RegisterStub(ip.String(), "tabellone", utils.Replicas, utils.AwsRegion)
	for len(replicaSet) != utils.Replicas {
		log.Println("Waiting for replicas to connect...")
		time.Sleep(60 * time.Second)
		replicaSet = cloud.RegisterStub(ip.String(), "tabellone", utils.Replicas, utils.AwsRegion)
	}

	log.Println("Replicas found: ", replicaSet)

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
