package main

import (
	"SDCC/client"
	_ "SDCC/client"
	"SDCC/utils"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"strings"

	"github.com/manifoldco/promptui"
)

var commands = []string{"get", "put", "del", "append", "nodes", "connect", "lnodes", "help", "closest", "disconnect"}

var conn *client.Connection
var connected = false
var grpcConn *grpc.ClientConn

func main() {
	validate := func(input string) error {

		parts := strings.Split(input, " ")
		command := parts[0]

		if !utils.Contains(commands, command) {
			return errors.New("invalid command")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Command",
		Validate: validate,
	}

	for {
		result, err := prompt.Run()

		parts := strings.Split(result, " ")
		command := parts[0]

		switch command {
		case "get":
			if len(parts) == 2 {
				if !connected {
					fmt.Println("Not connected to a client")
				} else {
					ret, err := conn.Get([]byte(parts[1]))
					if err != nil {
						fmt.Println("Error in get: ", err)
					} else {
						for _, slice := range ret {
							fmt.Println(string(slice))
						}
					}
				}
			} else {
				fmt.Println("Invalid number of parameters: get <key>")
			}
		case "put":
			if len(parts) >= 3 {
				if !connected {
					fmt.Println("Not connected to a client")
				} else {
					args := len(parts[2:])
					input := make([][]byte, args)

					for i, v := range parts[2:] {
						input[i] = []byte(v)
					}

					err := conn.Put([]byte(parts[1]), input)
					if err != nil {
						fmt.Println("Error in put: ", err)
					} else {
						fmt.Println("OK")
					}
				}
			} else {
				fmt.Println("Invalid number of parameters: put <key> <value>")
			}
		case "lnodes":
			if len(parts) == 2 {
				allNodes, err := client.GetAllNodesLocal(parts[1])
				if err != nil {
					fmt.Printf("Error: failed to contact %s", parts[1])
				} else {
					fmt.Println("Nodes: ", allNodes)
				}
			} else {
				fmt.Println("Invalid number of parameters: lnodes <address>")
			}

		case "closest":
			if len(parts) == 2 {
				allNodes, err := client.GetAllNodesLocal(parts[1])
				if err != nil {
					fmt.Printf("Error: failed to contact %s", parts[1])
				}

				node, err := client.GetClosestNode(allNodes)
				if err != nil {
					fmt.Println("Error: ", err)
				} else {
					fmt.Println("Closest node is ", node)
				}

			} else {
				fmt.Println("Invalid number of parameters: closest <address>")
			}

		case "connect":
			if len(parts) == 2 {
				conn, grpcConn, err = client.Connect(parts[1])
				if err != nil {
					fmt.Printf("Error: failed to contact %s", parts[1])
				} else {
					fmt.Println("Connection successful")
					connected = true
				}
			} else {
				fmt.Println("Invalid number of parameters: connect <address>")
			}
		case "disconnect":
			err := grpcConn.Close()
			if err != nil {
				fmt.Println("Error while closing")
			}

			connected = false
		}

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
	}
}
