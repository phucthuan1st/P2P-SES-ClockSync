package node

import (
	"fmt"
	"net"
	"os"
	"p2p-ses-clocksync/vectorclock"
	"sync"
)

type Node struct {
	ID              string // ID is own remote address as ip:port
	IP              string
	Port            string
	Mutex           sync.Mutex
	Listener        net.Listener        // own listener
	NodesConnection map[string]net.Conn // save connection of connected nodes
	OwnVectorClock  *vectorclock.VectorClock
	OtherNodeClock  map[string]*vectorclock.VectorClock
}

func NewNode(ip, port string) *Node {
	p := &Node{
		ID:              ip + ":" + port,
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
		p.AddNodeConnection(remoteAddr, conn)
	}
}

func (p *Node) AddNodeConnection(id string, conn net.Conn) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.NodesConnection[id] = conn
	fmt.Printf("[a] Node %s add Node %s to its handler\n", p.ID, id)
	go p.HandleNodeCommunication(id, conn)
}

func (p *Node) HandleNodeCommunication(id string, conn net.Conn) {
	defer func() {
		conn.Close()
		p.RemoveNodeConnection(id)
	}()
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			fmt.Printf("[x] Node %s disconnected from Node %s\n", p.ID, id)
			break
		}
		message := string(data[:n])
		fmt.Printf("<-- %s received from Node %s: %s\n", p.ID, id, message)
	}
}

func (p *Node) RemoveNodeConnection(id string) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	delete(p.NodesConnection, id)
}

func (p *Node) SendMessage(id, message string) {
	fmt.Printf("--> Node %s sending message to node %s: %s\n", p.ID, id, message)
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if conn, ok := p.NodesConnection[id]; ok {
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message to Node", id)
		}
	} else {
		fmt.Println("Node", id, "not found")
	}
}

func (p *Node) MakeNodeConnection(nodeIP, nodePort string) error {
	remoteAddr := nodeIP + ":" + nodePort

	// Establish a TCP connection to the remote peer
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		fmt.Printf("[x] Error connecting to Node %s at %s: %v\n", remoteAddr, remoteAddr, err)
		return err
	}

	// Add the connected peer to the list of peers
	fmt.Printf("[v] Node %s make a connection to Node %s\n", p.ID, remoteAddr)
	p.AddNodeConnection(remoteAddr, conn)

	return nil
}
