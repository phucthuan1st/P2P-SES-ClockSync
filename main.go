package main

import (
	"math/rand"
	"net"
	"p2p-ses-clocksync/node"
	"time"
)

func main() {

	// Create three nodes
	nodeA := node.NewNode("127.0.0.1", "9000")
	nodeB := node.NewNode("127.0.0.1", "9001")
	nodeC := node.NewNode("127.0.0.1", "9002")

	// Connect nodes to each other
	nodeA.MakeNodeConnection("127.0.0.1", "9001")
	nodeA.MakeNodeConnection("127.0.0.1", "9002")

	nodeB.MakeNodeConnection("127.0.0.1", "9002")

	// Simulate sending messages
	go func() {
		message := "Hello from A"
		for i := 0; i < 5; i++ {
			randomNodeID := getRandomNodeID(nodeA.NodesConnection)
			nodeA.SendMessage(randomNodeID, message)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		message := "Hello from B"

		for i := 0; i < 5; i++ {
			randomNodeID := getRandomNodeID(nodeB.NodesConnection)
			nodeB.SendMessage(randomNodeID, message)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		message := "Hello from C"

		for i := 0; i < 5; i++ {
			randomNodeID := getRandomNodeID(nodeC.NodesConnection)
			nodeC.SendMessage(randomNodeID, message)
			time.Sleep(time.Second)
		}
	}()

	// Keep the program running
	select {}
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
