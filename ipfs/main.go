package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	Port           int
	Rendezvous     string
	Seed           int64
	DiscoveryPeers addrList
	TestMode       bool
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

func main() {
	config := Config{}

	flag.StringVar(&config.Rendezvous, "rendezvous", "sdcc", "")
	flag.Int64Var(&config.Seed, "seed", 0, "Seed value for generating a PeerID, 0 is random")
	flag.Var(&config.DiscoveryPeers, "peer", "Peer multiaddress for peer discovery")
	flag.IntVar(&config.Port, "port", 0, "")
	flag.BoolVar(&config.TestMode, "test-mode", false, "")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	h, err := NewHost(ctx, config.Seed, config.Port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Host ID: %s", h.ID().Pretty())
	log.Printf("Connect to me on:")
	for _, addr := range h.Addrs() {
		log.Printf("  %s/p2p/%s", addr, h.ID().Pretty())
	}

	dht, err := NewDHT(ctx, h, config.DiscoveryPeers)
	dht.Validator = NullValidator{}
	if err != nil {
		log.Fatal(err)
	}

	if config.TestMode {
		/*err := dht.PutValue(ctx, "ciao", []byte("ciao"))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("OK") */
		value, err := dht.SearchValue(ctx, "ciao")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(string(<-value))
		return
	}

	// go Discover(ctx, h, dht, config.Rendezvous)
	if len(config.DiscoveryPeers) != 0 {
		ids, err := dht.GetClosestPeers(ctx, peer.Encode(h.ID()))
		for _, id := range ids {
			log.Printf(peer.Encode(id))
		}
		if err != nil {
			log.Print(err)
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()
		for {
			 <-ticker.C
			 log.Printf("Peers")
			 for _, id := range dht.RoutingTable().ListPeers() {
				 log.Printf(peer.Encode(id))
			 }
			 log.Printf("Table")
			 dht.RoutingTable().Print()
		}
	}()

	run(h, cancel)
}

func run(h host.Host, cancel func()) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Printf("\rExiting...\n")

	cancel()

	if err := h.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
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
