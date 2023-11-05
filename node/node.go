package node

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"p2p-ses-clocksync/message"
	"p2p-ses-clocksync/vectorclock"
	"sync"
)

func enqueue(queue []message.Message, element message.Message) []message.Message {
	queue = append(queue, element) // Simply append to enqueue.
	return queue
}

func dequeue(queue []message.Message) (message.Message, []message.Message) {
	element := queue[0] // The first element is the one to be dequeued.
	if len(queue) == 1 {
		var tmp = []message.Message{}
		return element, tmp

	}

	return element, queue[1:] // Slice off the element once it is dequeued.
}

type Node struct {
	ID              string // ID is own remote address as ip:port
	IP              string
	Port            string
	Mutex           sync.RWMutex
	Listener        net.Listener        // own listener
	NodesConnection map[string]net.Conn // save connection of connected nodes
	OwnVectorClock  *vectorclock.VectorClock
	OtherNodeClock  map[string]*vectorclock.VectorClock
	MessageBuffer   []message.Message
}

func NewNode(id, ip, port string) *Node {

	p := &Node{
		ID:              id,
		IP:              ip,
		Port:            port,
		NodesConnection: make(map[string]net.Conn),
		OwnVectorClock:  vectorclock.NewVectorClock(),
		OtherNodeClock:  make(map[string]*vectorclock.VectorClock),
		MessageBuffer:   make([]message.Message, 0),
	}
	p.StartListener()
	return p
}

func (node *Node) StartListener() {
	listenAddr := node.IP + ":" + node.Port
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("[x] Error starting listener:", err)
		os.Exit(1)
	}
	node.Listener = listener
	log.Printf("[o] Node %s start listening on %s\n", node.ID, listenAddr)
	go node.AcceptConnections()
}

func (node *Node) AcceptConnections() {
	for {
		conn, err := node.Listener.Accept()
		if err != nil {
			fmt.Println("[x] Error accepting connection:", err)
			continue
		}
		remoteAddr := conn.RemoteAddr().String()

		// Read the first message (name) from the incoming node
		nameBuffer := make([]byte, 128) // Adjust the buffer size as needed
		n, err := conn.Read(nameBuffer)
		if err != nil {
			log.Printf("[x] Error reading name from Node %s: %v\n", remoteAddr, err)
			continue
		}
		incomingName := string(nameBuffer[:n])

		// Add the connected peer to the list of peers with the incoming name
		log.Printf("[a] Node %s accepted connection from Node %s\n", node.ID, incomingName)
		node.AddNodeConnection(incomingName, conn)
	}
}

func (node *Node) AddNodeConnection(id string, conn net.Conn) {
	node.Mutex.Lock()
	node.NodesConnection[id] = conn
	node.Mutex.Unlock()

	go node.HandleNodeCommunication(id, conn)
}

func (node *Node) RemoveNodeConnection(id string) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()
	delete(node.NodesConnection, id)
}

func (node *Node) MakeNodeConnection(nodeID, nodeIP, nodePort string) error {
	remoteAddr := nodeIP + ":" + nodePort

	// Establish a TCP connection to the remote peer
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("[x] Error connecting to Node %s at %s: %v\n", nodeID, remoteAddr, err)
		return err
	}

	// Add the connected peer to the list of peers
	log.Printf("[m] Node %s make a connection to Node %s\n", node.ID, nodeID)
	node.AddNodeConnection(nodeID, conn)

	// Send a message with the name of the node
	message := node.ID
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Printf("[x] Error sending message to Node %s: %v\n", nodeID, err)
		return err
	}

	return nil
}

func (node *Node) HandleNodeCommunication(id string, conn net.Conn) {
	defer func() {
		conn.Close()
		node.RemoveNodeConnection(id)
	}()

	decoder := json.NewDecoder(conn)
	for {
		var msg message.Message

		// Read one line (up to the newline delimiter)
		if err := decoder.Decode(&msg); err != nil {
			log.Printf("[x] Node %s disconnected from Node %s\n", node.ID, id)
			break
		}

		// Handle the received message
		go node.handleReceivedMessage(id, msg) // Use a goroutine to handle each message concurrently
	}
}

func (node *Node) handleReceivedMessage(nodeSrcId string, msg message.Message) {

	// TODO: if there is no line in payloads, or no exist line in payloads match the ID
	// then merge max value for all clock entries, and deliver the message
	if index, isContained := containsName(msg.Payloads, node.ID); !isContained {
		node.DeliverMessage(msg, nil)

		// TODO: check the buffer if there is any message that can be delivered, if yes, deliver it
		var bufferedMsg message.Message
		node.Mutex.Lock()
		for len(node.MessageBuffer) > 0 {
			bufferedMsg, node.MessageBuffer = dequeue(node.MessageBuffer)
			go node.handleReceivedMessage(bufferedMsg.Source, bufferedMsg)
		}
		node.Mutex.Unlock()
	} else {
		payload := msg.Payloads[index]
		payloadClock := vectorclock.NewVectorClock()
		payloadClock.SetClock(payload.Clock)

		// TODO: if t <= local clock then deliver the message
		if payloadClock.Compare(node.OwnVectorClock) == -1 {
			node.DeliverMessage(msg, &payload)

			// TODO: check the buffer if there is any message that can be delivered, if yes, deliver it
			var bufferedMsg message.Message
			node.Mutex.Lock()
			for len(node.MessageBuffer) > 0 {
				bufferedMsg, node.MessageBuffer = dequeue(node.MessageBuffer)
				go node.handleReceivedMessage(bufferedMsg.Source, bufferedMsg)
			}
			node.Mutex.Unlock()
		}

		// TODO: buffer the message
		node.BufferMessage(msg, payload)
	}
}

func (node *Node) DeliverMessage(msg message.Message, payload *message.Payload) {
	node.Mutex.Lock()

	log.Printf("[+] Message: [%s] sent from Node %s to Node %s\n", msg.Content, msg.Source, node.ID)
	log.Printf("    Content: %s\n", msg.Content)
	log.Printf("    Status: Delivered at %v\n", msg.Timestamp)

	if payload != nil {
		log.Printf("    Cause: Message from %s at %v is delivered\n", payload.Name, payload.Clock)
	}

	// Log the changes in clock
	for name, clock := range node.OtherNodeClock {
		log.Printf("    Clock for Node %s (before merge): %v\n", name, clock.GetClock())
	}

	/*
		Merge V_M (in message) with V_P2 as follows.
			If (P,t) is not there in V_P2, merge.
			If (P,t) is present in V_P2, t is updated with max(t[i] in Vm, t[i] in V_P2). {Component-wise maximum}.
	*/
	for _, payload := range msg.Payloads {
		clock := vectorclock.NewVectorClock()
		clock.SetClock(payload.Clock)

		if _, ok := node.OtherNodeClock[payload.Name]; !ok {
			node.OtherNodeClock[payload.Name] = clock
		} else {
			clockEntries := node.OtherNodeClock[payload.Name].Clone().GetClock()
			node.OtherNodeClock[payload.Name] = clock.Merge(clockEntries)
		}

		// Log the changes in clock after the merge
		log.Printf("    Clock for Node %s (after merge): %v\n", payload.Name, node.OtherNodeClock[payload.Name].GetClock())
	}

	// Update site P2’s local, logical clock.
	node.OwnVectorClock.Increment(node.ID)
	node.OwnVectorClock = node.OwnVectorClock.Merge(msg.Timestamp)

	log.Printf("    Updated Clock for Node %s: %v\n", node.ID, node.OwnVectorClock.GetClock())
	node.Mutex.Unlock()
}

func (node *Node) BufferMessage(msg message.Message, cause message.Payload) {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	log.Printf("[-] Message: [%s] sent from Node %s to Node %s\n", msg.Content, msg.Source, node.ID)
	log.Printf("    Content: %s\n", msg.Content)
	log.Printf("    Status: Buffered\n")
	log.Printf("    Cause: wait for delivery to %s at %v\n", cause.Name, cause.Clock)
	log.Printf("    Timestamp: %v\n", msg.Timestamp)
	log.Printf("    Local time: %v\n", node.OwnVectorClock.GetClock())

	node.MessageBuffer = enqueue(node.MessageBuffer, msg)
}

func containsName(payloads []message.Payload, name string) (int, bool) {
	for index, payload := range payloads {
		if payload.Name == name {
			return index, true
		}
	}
	return -1, false
}

func (node *Node) SendMessage(destNodeID, content string) error {
	node.Mutex.Lock()
	defer node.Mutex.Unlock()

	// Check if the destination node exists in the connections
	conn, ok := node.NodesConnection[destNodeID]
	if !ok {
		return fmt.Errorf("Node %s not found", destNodeID)
	}

	// Increment the own vector clock
	node.OwnVectorClock.Increment(node.ID)

	// Create payloads from other node clocks
	payloads := make([]message.Payload, 0)
	for name, clock := range node.OtherNodeClock {
		payloads = append(payloads, message.Payload{
			Name:  name,
			Clock: clock.GetClock(),
		})
	}

	// Create the message
	msg := message.Message{
		Source:    node.ID,
		Dest:      destNodeID,
		Content:   content,
		Timestamp: node.OwnVectorClock.Clone().GetClock(),
		Payloads:  payloads,
	}

	// Serialize the message
	serializedMessage, err := msg.Serialize()
	if err != nil {
		return fmt.Errorf("Serialization failed: %v", err)
	}

	serializedMessage = serializedMessage + "\n"

	// Send the message with content M, timestamp tm, and other nodes' clock VPs
	_, err = conn.Write([]byte(serializedMessage))
	if err != nil {
		return fmt.Errorf("Failed to send message to Node %s: %v", destNodeID, err)
	}

	// Add (Pj, tm) to other nodes' clock VP, rewrite if exists
	node.OtherNodeClock[destNodeID] = node.OwnVectorClock.Clone()

	log.Printf("[+] Message: [%s] sent from Node %s to Node %s at %v\n", content, node.ID, destNodeID, msg.Timestamp)

	return nil
}
