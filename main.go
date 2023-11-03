package main

import (
	"fmt"
	"math/rand"
	"net"
	"p2p-ses-clocksync/node"
)

func main() {

	// TODO: Create three nodes
	nodeA := node.NewNode("A", "127.0.0.1", "9000")
	nodeB := node.NewNode("B", "127.0.0.1", "9001")
	nodeC := node.NewNode("C", "127.0.0.1", "9002")

	// TODO: Connect nodes to each other
	nodeA.MakeNodeConnection("B", "127.0.0.1", "9001")
	nodeA.MakeNodeConnection("C", "127.0.0.1", "9002")

	nodeB.MakeNodeConnection("C", "127.0.0.1", "9002")

	// TODO: After complete connect each node, start to simulate message sending and clock synchronization
	go func() {
		message := "Hello from A"
		for i := 0; i < 5; i++ {
			randomNodeID := getRandomNodeID(nodeA.NodesConnection)
			nodeA.SendMessage(randomNodeID, fmt.Sprintf("%d - %s", i, message))
			//time.Sleep(time.Second)
		}
	}()

	go func() {
		message := "Hello from B"

		for i := 0; i < 5; i++ {
			randomNodeID := getRandomNodeID(nodeB.NodesConnection)
			nodeB.SendMessage(randomNodeID, fmt.Sprintf("%d - %s", i, message))
			//time.Sleep(time.Second)
		}
	}()

	go func() {
		message := "Hello from C"

		for i := 0; i < 5; i++ {
			randomNodeID := getRandomNodeID(nodeC.NodesConnection)
			nodeC.SendMessage(randomNodeID, fmt.Sprintf("%d - %s", i, message))
			//time.Sleep(time.Second)
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
