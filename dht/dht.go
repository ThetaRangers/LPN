package dht

import (
	"SDCC/utils"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

const(
	valueChar = "0"
	clusterChar = "1"
	offloadChar = "2"
	localChar = "3"
)

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

type KDht struct {
	kDht *dht.IpfsDHT
}

func (k *KDht) GetValue(ctx context.Context, key string) (string, bool, error) {
	dhtKey := valueChar + key
	value, err := k.kDht.GetValue(ctx, dhtKey)
	if err != nil {
		return "", false, err
	} else if string(value[0]) == offloadChar {
		return string(value[1:]), true, nil
	} else {
		return string(value[1:]), false, nil
	}
}

func (k *KDht) PutValue(ctx context.Context, key string, value string, offloaded bool) error {
	dhtKey := valueChar + key
	var dhtValue = make([]byte, 1)
	if offloaded {
		dhtValue[0] = offloadChar[0]
	} else {
		dhtValue[0] = localChar[0]
	}
	dhtValue = append(dhtValue, []byte(value)...)
	return k.kDht.PutValue(ctx, dhtKey, dhtValue)
}

func (k *KDht) GetCluster(ctx context.Context, key string) ([]string, error) {
	dhtKey := clusterChar + key
	value, err := k.kDht.GetValue(ctx, dhtKey)
	if err != nil {
		return nil, err
	} else {
		var target []string
		err := json.Unmarshal(value, &target)
		if err != nil {
			return nil, err
		}
		return target, nil
	}
}

func (k *KDht) PutCluster(ctx context.Context, key string, value *utils.ClusterRoutine) error {
	dhtKey := clusterChar + key
	cluster := value.GetAll()
	dhtValue, err := json.Marshal(cluster)
	if err != nil {
		return err
	}
	return k.kDht.PutValue(ctx, dhtKey, dhtValue)
}

func NewKDht(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (*KDht, error) {
	var options []dht.Option

	options = append(options, dht.Mode(dht.ModeServer))
	options = append(options, dht.ProtocolPrefix("/sdcc"))

	ipfsDHT, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = ipfsDHT.Bootstrap(ctx); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeers {
		peerInfo, _ := peer.AddrInfoFromP2pAddr(peerAddr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerInfo); err != nil {
				log.Printf("Error while connecting to node %q: %-v", peerInfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %q", *peerInfo)
			}
		}()
	}
	wg.Wait()

	ipfsDHT.Validator = NullValidator{}
	return &KDht{kDht: ipfsDHT}, nil
}
