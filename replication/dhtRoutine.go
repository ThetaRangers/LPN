package replication

import (
	"SDCC/dht"
	"SDCC/utils"
	"context"
	"log"
)

type DhtRoutine struct {
	alive   chan *dht.KDht
	op      chan struct{}
	kDht    *dht.KDht
	cluster *utils.ClusterRoutine
	ip      string
}

func (d DhtRoutine) run() {
	d.kDht = <-d.alive
	close(d.alive)
	for {
		<-d.op
		err := d.kDht.PutCluster(context.Background(), d.ip, d.cluster)
		if err != nil {
			log.Println(err)
		}
	}
}

func (d DhtRoutine) SetDht(kDht *dht.KDht) {
	d.alive <- kDht
}

func (d DhtRoutine) UpdateCluster() {
	var foo struct{}
	d.op <- foo
}

func NewDhtRoutine(ip string, cluster *utils.ClusterRoutine) DhtRoutine {
	ret := DhtRoutine{alive: make(chan *dht.KDht), op: make(chan struct{}, 200), cluster: cluster, ip: ip}
	go ret.run()
	return ret
}
