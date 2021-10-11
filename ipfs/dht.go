package ipfs

import (
	"context"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

func NewDHT(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (*dht.IpfsDHT, error) {
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

	return ipfsDHT, nil
}
