package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ThetaRangers/Kademlia/dht"
	"github.com/urfave/cli"
)

func prepareArgs() *cli.App {
	cli.AppHelpTemplate = `NAME:
	{{.Name}} - {{.Usage}}
USAGE:
	{{if .VisibleFlags}}{{.HelpName}} [options]{{end}}
	{{if len .Authors}}
AUTHOR:
	{{range .Authors}}{{ . }}{{end}}
	{{end}}{{if .Commands}}
VERSION:
	{{.Version}}
OPTIONS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
	{{.Copyright}}
	{{end}}{{if .Version}}
	{{end}}`

	cli.VersionFlag = cli.BoolFlag{
		Name:  "V, version",
		Usage: "Print version",
	}

	cli.HelpFlag = cli.BoolFlag{
		Name:  "h, help",
		Usage: "Print help",
	}

	app := cli.NewApp()

	app.Name = "DHT"
	app.Version = "0.2.0"
	app.Compiled = time.Now()

	app.Usage = "Experimental Distributed Hash Table"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "c, connect",
			Usage: "Connect to bootstrap node ip:port",
		},
		cli.StringFlag{
			Name:  "l, listen",
			Usage: "Listening address and port",
			Value: ":3000",
		},
		cli.BoolFlag{
			Name:  "i, interactif",
			Usage: "Interactif",
		},
		cli.BoolFlag{
			Name:  "s, store",
			Usage: "Store from Stdin",
		},
		cli.StringFlag{
			Name:  "S, store-at",
			Usage: "Same as '-s' but store at given `key`",
		},
		cli.StringFlag{
			Name:  "f, fetch",
			Usage: "Fetch `hash` and prints to Stdout",
		},
		cli.StringFlag{
			Name:  "F, fetch-at",
			Usage: "Same as '-f' but fetch from given `key`",
		},
		cli.IntFlag{
			Name:  "n, network",
			Value: 0,
			Usage: "Spawn X new `nodes` in a network.",
		},
		cli.IntFlag{
			Name:  "v, verbose",
			Value: 3,
			Usage: "Verbose `level`, 0 for CRITICAL and 5 for DEBUG",
		},
	}

	app.UsageText = "dht [options]"

	return app
}

func parseArgs(done func(dht.DhtOptions, *cli.Context)) {
	app := prepareArgs()

	app.Action = func(c *cli.Context) error {
		options := dht.DhtOptions{
			ListenAddr:    c.String("l"),
			BootstrapAddr: c.String("c"),
			Verbose:       c.Int("v"),
			Stats:         c.Bool("s"),
			Interactif:    c.Bool("i"),
			Cluster:       c.Int("n"),
			// OnStore:       func(dht.Packet) interface{} {},
		}

		if options.Cluster > 0 {
			options.Stats = false
			options.Interactif = false
		}

		if options.Interactif {
			options.Stats = false
		}

		done(options, c)

		return nil
	}

	app.Run(os.Args)
}

func main() {
	parseArgs(func(options dht.DhtOptions, c *cli.Context) {
		if options.Cluster > 0 {
			cluster(options)
		} else {
			node := startOne(options)

			if c.Bool("s") {
				storeFromStdin(node)
			} else if len(c.String("S")) > 0 {
				storeAt(node, c.String("S"))
			} else if len(c.String("f")) > 0 {
				fetchFromHash(node, c.String("f"))
			} else if len(c.String("F")) > 0 {
				fetchAt(node, c.String("F"))
			} else {
				node.Wait()
			}

			listenExitSignals(node)
		}
	})
}

func storeFromStdin(node *dht.Dht) {
	res, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		node.Logger().Critical("Cannot store", err)

		return
	}

	hash, nb, err := node.Store(res)

	if err != nil {
		node.Logger().Critical("Cannot store", err)

		return
	}

	if nb == 0 {
		node.Logger().Critical("Cannot store, no nodes found")

		return
	}

	fmt.Println(hex.EncodeToString(hash))
}

func storeAt(node *dht.Dht, hashStr string) {
	res, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		node.Logger().Critical("Cannot store", err)

		return
	}

	hash := dht.NewHash([]byte(hashStr))

	_, nb, err := node.StoreAt(hash, res)

	if err != nil {
		node.Logger().Critical("Cannot store", err)

		return
	}

	if nb == 0 {
		node.Logger().Critical("Cannot store, no nodes found")

		return
	}

	fmt.Println(hex.EncodeToString(hash))
}

func fetchFromHash(node *dht.Dht, hashStr string) {
	hash, err := hex.DecodeString(hashStr)

	if err != nil {
		node.Logger().Critical("Cannot fetch", err)

		return
	}

	// time.Sleep(time.Second)

	b, err := node.Fetch(hash)

	if err != nil {
		node.Logger().Critical("Cannot fetch", err)

		return
	}

	fmt.Print(string(b))
}

func fetchAt(node *dht.Dht, hashStr string) {
	hash := dht.NewHash([]byte(hashStr))

	// time.Sleep(time.Second)

	b, err := node.Fetch(hash)

	if err != nil {
		node.Logger().Critical("Cannot fetch", err)

		return
	}

	fmt.Print(string(b))
}

func listenExitSignals(client *dht.Dht) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs

		exitProperly(client)

		os.Exit(0)
	}()
}

func exitProperly(client *dht.Dht) {
	client.Stop()
}

func cluster(options dht.DhtOptions) {
	network := []*dht.Dht{}
	i := 0
	if len(options.BootstrapAddr) == 0 {
		client := startOne(options)

		network = append(network, client)

		options.BootstrapAddr = options.ListenAddr

		i++
	}

	for ; i < options.Cluster; i++ {
		options2 := options

		addrPort := strings.Split(options.ListenAddr, ":")

		addr := addrPort[0]

		port, _ := strconv.Atoi(addrPort[1])

		options2.ListenAddr = addr + ":" + strconv.Itoa(port+i)

		client := startOne(options2)

		network = append(network, client)
	}

	for {
		time.Sleep(time.Second)
	}
}

func startOne(options dht.DhtOptions) *dht.Dht {
	client := dht.New(options)

	if err := client.Start(); err != nil {
		client.Logger().Critical(err)
	}

	return client
}
