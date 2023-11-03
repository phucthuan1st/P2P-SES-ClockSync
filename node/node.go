package node

import (
	"fmt"
	"net"
	"os"
	"p2p-ses-clocksync/vectorclock"
	"strings"
	"sync"
)

type Node struct {
	ID              string // ID is own remote address as ip:port
	IP              string
	Port            string
	Mutex           sync.RWMutex
	Listener        net.Listener        // own listener
	NodesConnection map[string]net.Conn // save connection of connected nodes
	OwnVectorClock  *vectorclock.VectorClock
	OtherNodeClock  map[string]*vectorclock.VectorClock
}

func NewNode(id, ip, port string) *Node {
	p := &Node{
		ID:              id,
		IP:              ip,
		Port:            port,
		NodesConnection: make(map[string]net.Conn),
		OwnVectorClock:  vectorclock.NewVectorClock(),
		OtherNodeClock:  make(map[string]*vectorclock.VectorClock),
	}
	p.StartListener()
	return p
}

func (p *Node) StartListener() {
	listenAddr := p.IP + ":" + p.Port
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("[x] Error starting listener:", err)
		os.Exit(1)
	}
	p.Listener = listener
	fmt.Printf("[o] Node %s start listening on %s\n", p.ID, listenAddr)
	go p.AcceptConnections()
}

func (p *Node) AcceptConnections() {
	for {
		conn, err := p.Listener.Accept()
		if err != nil {
			fmt.Println("[x] Error accepting connection:", err)
			continue
		}
		remoteAddr := conn.RemoteAddr().String()

		// Read the first message (name) from the incoming node
		nameBuffer := make([]byte, 128) // Adjust the buffer size as needed
		n, err := conn.Read(nameBuffer)
		if err != nil {
			fmt.Printf("[x] Error reading name from Node %s: %v\n", remoteAddr, err)
			continue
		}
		incomingName := string(nameBuffer[:n])

		// Add the connected peer to the list of peers with the incoming name
		fmt.Printf("[a] Node %s accepted connection from Node %s\n", p.ID, incomingName)
		p.AddNodeConnection(incomingName, conn)
	}
}

func (p *Node) AddNodeConnection(id string, conn net.Conn) {
	p.Mutex.Lock()
	p.NodesConnection[id] = conn
	p.Mutex.Unlock()

	go p.HandleNodeCommunication(id, conn)
}

func (p *Node) HandleNodeCommunication(id string, conn net.Conn) {
	defer func() {
		conn.Close()
		p.RemoveNodeConnection(id)
	}()
	for {
		data := make([]byte, 128)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Printf("[x] Node %s disconnected from Node %s\n", p.ID, id)
			break
		}

		// TODO: check the clock to decide whether buffer or delivery
		dataString := string(data[:n])
		messages := strings.Split(dataString, "|")

		for _, message := range messages {
			if message != "" {
				fmt.Printf("<-- %s received a message from Node %s: %s\n", p.ID, id, message)
			}
		}
	}
}

func (p *Node) RemoveNodeConnection(id string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	delete(p.NodesConnection, id)
}

func (p *Node) SendMessage(id, message string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()

	if conn, ok := p.NodesConnection[id]; ok {
		// TODO: send message along with the other clock
		// increment own process clock counter
		p.OwnVectorClock.Increment(p.ID)
		fmt.Printf("--> Node %s sending message to node %s: %s\n", p.ID, id, message)
		_, err := conn.Write([]byte(message + "|"))

		if err != nil {
			fmt.Println("Error sending message to Node", id)
		}
	} else {
		fmt.Println("Node", id, "not found")
	}
}

func (p *Node) MakeNodeConnection(nodeID, nodeIP, nodePort string) error {
	remoteAddr := nodeIP + ":" + nodePort

	// Establish a TCP connection to the remote peer
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("[x] Error connecting to Node %s at %s: %v\n", nodeID, remoteAddr, err)
		return err
	}

	// Add the connected peer to the list of peers
	fmt.Printf("[m] Node %s make a connection to Node %s\n", p.ID, nodeID)
	p.AddNodeConnection(nodeID, conn)

	// Send a message with the name of the node
	message := p.ID
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("[x] Error sending message to Node %s: %v\n", nodeID, err)
		return err
	}

	return nil
}
