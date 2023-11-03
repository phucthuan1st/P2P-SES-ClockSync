package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"p2p-ses-clocksync/node"
)

type NodeConfig struct {
	NodeID    string
	NodeIP    string
	NodePort  int
	ConnectTo []struct {
		NodeID   string
		NodeIP   string
		NodePort int
	}
}

func main() {
	configFile := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configFile == "" {
		flag.Usage()
		fmt.Println("Error: --config is required.")
		os.Exit(1)
	}

	configData, err := os.ReadFile(*configFile)
	if err != nil {
		fmt.Printf("Error reading configuration file: %v\n", err)
		os.Exit(1)
	}

	var nodeConfig NodeConfig
	if err := json.Unmarshal(configData, &nodeConfig); err != nil {
		fmt.Printf("Error parsing configuration file: %v\n", err)
		os.Exit(1)
	}

	current := node.NewNode(nodeConfig.NodeID, nodeConfig.NodeIP, fmt.Sprint(nodeConfig.NodePort))

	// Menu for choosing connect a new node, or start simulating
	for {
		// Display the menu
		fmt.Println("Choose an option:")
		fmt.Println("1. Connect to a new node")
		fmt.Println("2. Start simulating")
		fmt.Println("3. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			// Implement connecting to a new node
			for _, connection := range nodeConfig.ConnectTo {
				current.MakeNodeConnection(connection.NodeID, connection.NodeIP, fmt.Sprint(connection.NodePort))
			}
		case 2:
			go simulateMessages(current)
		case 3:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please enter a valid option (1, 2, or 3).")
		}
	}
}

func simulateMessages(node *node.Node) {
	message := fmt.Sprintf("Hello from %s", node.ID)
	for i := 0; i < 5; i++ {
		randomNodeID := getRandomNodeID(node.NodesConnection)
		node.SendMessage(randomNodeID, fmt.Sprintf("%d - %s", i, message))
	}
}

func getRandomNodeID(connections map[string]net.Conn) string {
	var nodeIDs []string
	for nodeID := range connections {
		nodeIDs = append(nodeIDs, nodeID)
	}
	if len(nodeIDs) > 0 {
		randomIndex := rand.Intn(len(nodeIDs))
		return nodeIDs[randomIndex]
	}
	return ""
}
