package server

import (
	"context"
	"fmt"
	"go_server/utils"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type LeaderNode struct {
	Uuid    string
	Address string
}

type Node struct {
	// Default struct items
	uuid         string
	nodeType     utils.Entity
	address      string
	leader       LeaderNode
	nodes        map[string]string     // uuid -> address
	nodeServices map[Service]*[]string // service -> []uuid
	serviceList  []Service             // List of services this node manages
	// For maintaining concurrency
	mtx sync.RWMutex
	// For closing down and cleanup
	cancel  context.CancelFunc
	context context.Context
	// For console management
	outputChannel chan string
}

func HandleNode(local string, remote string) {
	fmt.Println("Generating new node")
	// For cleanup
	context, cancel := context.WithCancel(context.Background())

	// Create the node
	node := Node{
		uuid:          strings.ReplaceAll(uuid.NewString(), "-", ""),
		nodeType:      utils.Node,
		address:       local,
		leader:        LeaderNode{},
		nodes:         make(map[string]string),
		nodeServices:  make(map[Service]*[]string),
		cancel:        cancel,
		context:       context,
		outputChannel: make(chan string, 100),
	}
	node.initServices()
	node.nodes[node.uuid] = node.address

	// Listen for incoming
	go node.listen()

	// Set up Console
	var wait sync.WaitGroup
	wait.Add(1)
	go node.inputHandler(&wait)

	// Initialization
	if remote == "" {
		node.nodeType = utils.Leader
		node.leader.Uuid = node.uuid
		node.leader.Address = node.address
	} else {
		// Initialize the node with its services
		node.initServices()

		// Register with the leader
		node.outputChannel <- fmt.Sprintf("Attempting to register with Leader at: %s", remote)
		registerRequest, err := RegisterRequest(&node)
		if err != nil {
			fmt.Printf("Could not generate registration request: %v", err)
			node.shutdown()
		}

		// Connect with leader
		conn, err := net.Dial(string(utils.TCP), remote)
		if err != nil {
			node.outputChannel <- fmt.Sprintf("Could not generate a connection with leader at %s on %s: %v", remote, node.address, err)
			node.shutdown()
		}
		// Send the request
		conn.Write(registerRequest)
		conn.Close()
	}

	node.mtx.RLock()
	if node.nodeType != utils.Leader {
		node.printServiceList()
	}
	node.mtx.RUnlock()

	wait.Wait()
	node.shutdown()
}

func (node *Node) listen() {
	node.outputChannel <- "Initiating listener..."
	// Create the listener
	listener, err := net.Listen(string(utils.TCP), node.address)
	if err != nil {
		node.outputChannel <- fmt.Sprintf("Failed to listen on %s: %v\n", node.address, err)
		node.shutdown()
	}
	defer listener.Close()

	node.outputChannel <- fmt.Sprintf("Listening on %s\n", node.address)
	go node.outputHandler()

	// Listen for incoming connections
	for {
		select {
		case <-node.context.Done():
			node.outputChannel <- "Stopping listener"
			return

		default:
			connection, err := listener.Accept()
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Failed to accept connection: %v\n", err)
				continue
			}
			go node.handleConnection(connection)
		}
	}
}

func (node *Node) sendMessage(message []byte) {
	for id, address := range node.nodes {
		if id == node.uuid {
			continue
		}

		conn, err := net.Dial(string(utils.TCP), address)
		if err != nil {
			node.outputChannel <- fmt.Sprintf("Could not contact node %s :: %s", id, address)
			break
		}
		defer conn.Close()
		conn.Write(message)
	}
}

func (node *Node) shutdown() {
	// If current leader, ensure the other nodes are updated first
	if node.uuid == node.leader.Uuid && len(node.nodes) > 0 {
		// Change the leader
		for k, v := range node.nodes {
			if k != node.uuid {
				node.leader.Uuid = k
				node.leader.Address = v
				node.nodeType = utils.Node
				break
			}
		}
		// Remove from the nodes list
		delete(node.nodes, node.uuid)
		node.removeServiceNode(node.leader.Uuid)

		// Update the other nodes
		update, err := UpdateNodeListRequest(node)
		if err != nil {
			fmt.Println("No moar leadaaa")
		}
		node.sendMessage(update)
	}

	shutdown, err := ShutdownRequest(node)
	if err == nil {
		node.sendMessage(shutdown)
	} else {
		fmt.Println("Failed to prepare shutdown message")
	}

	close(node.outputChannel)
	fmt.Println("Exiting node")
	node.cancel()
}
