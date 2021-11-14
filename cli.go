package main

import (
	"SDCC/client"
	"fmt"
	"github.com/c-bata/go-prompt"
	"google.golang.org/grpc"
	"strings"
)

func completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		{Text: "get", Description: "Get value"},
		{Text: "put", Description: "Put key and value"},
		{Text: "connect", Description: "Connect to a node"},
		{Text: "disconnect", Description: "Disconnect from the node"},
		{Text: "quit", Description: "Quit from the cli"},
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

var history []string
var conn *client.Connection
var connected = false
var grpcConn *grpc.ClientConn

func main() {
	var err error

	for {

		t := prompt.Input(">>> ", completer,
			prompt.OptionTitle("lpn-promp"),
			prompt.OptionHistory(history),
			prompt.OptionPrefixTextColor(prompt.Green),
			prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
			prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
			prompt.OptionSuggestionBGColor(prompt.DarkGray))

		history = append(history, t)

		parts := strings.Split(t, " ")
		command := parts[0]

		switch command {
		case "quit":
			return
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

	}
}
