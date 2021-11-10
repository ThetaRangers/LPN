package main

import (
	"fmt"
	"github.com/hexablock/vivaldi"
	"labix.org/v1/vclock"
	"sync"
	"time"
)

func main() {
	conf := vivaldi.DefaultConfig()
	conf.Dimensionality = 2

	conf2 := vivaldi.DefaultConfig()
	conf2.Dimensionality = 2

	c, err := vivaldi.NewClient(conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.GetCoordinate())

	c2, err := vivaldi.NewClient(conf2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Default conf", conf)
	coor, err := c2.Update("node", c2.GetCoordinate(), 2*time.Second)
	fmt.Println("C1 coordinates: ", c.GetCoordinate())

	c2.SetCoordinate(coor)
	fmt.Println("C2 coordinates: ", c2.GetCoordinate())

	fmt.Println(c.DistanceTo(coor))

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		wg.Done()
	}()

	go func() {
		time.Sleep(2 * time.Second)
		wg.Done()
	}()

	time.Sleep(3 * time.Second)
	vc1 := vclock.New()
	vc1.Update("A", 1)

	vc2 := vc1.Copy()
	vc2.Update("B", 0)

	fmt.Println(vc2.Compare(vc1, vclock.Ancestor))   // true
	fmt.Println(vc1.Compare(vc2, vclock.Descendant)) // true

	vc1.Update("C", 5)

	fmt.Println(vc1.Compare(vc2, vclock.Descendant)) // false
	fmt.Println(vc1.Compare(vc2, vclock.Concurrent)) // true

	vc2.Merge(vc1)

	fmt.Println(vc1.Compare(vc2, vclock.Descendant)) // true

	data := vc2.Bytes()
	fmt.Printf("%#v\n", string(data))
}
