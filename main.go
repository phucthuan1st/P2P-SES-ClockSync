package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"p2p-ses-clocksync/node"
)

func main() {
	id := flag.String("id", "A", "Node ID")
	host := flag.String("host", "127.0.0.1", "Host IP")
	var port string
	flag.StringVar(&port, "port", "9000", "Port to listen")
	flag.Parse()

	if port == "" {
		flag.Usage()
		fmt.Println("Error: --port is required.")
		os.Exit(1)
	}

	current := node.NewNode(*id, *host, port)

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
			var nodeID, nodeIP, nodePort string

			fmt.Print("Enter the new node's ID: ")
			_, err := fmt.Scan(&nodeID)
			if err != nil {
				fmt.Println("Error reading input.")
				continue
			}

			fmt.Print("Enter the new node's IP: ")
			_, err = fmt.Scan(&nodeIP)
			if err != nil {
				fmt.Println("Error reading input.")
				continue
			}

			fmt.Print("Enter the new node's Port: ")
			_, err = fmt.Scan(&nodePort)
			if err != nil {
				fmt.Println("Error reading input.")
				continue
			}
			// Use nodeID, nodeIP, and nodePort to connect to the new node
			current.MakeNodeConnection(nodeID, nodeIP, nodePort)
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
