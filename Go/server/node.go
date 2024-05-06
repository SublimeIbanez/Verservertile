package server

import (
	"bufio"
	"context"
	"fmt"
	"go_server/utils"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type LeaderNode struct {
	uuid          string
	leaderAddress string
}

type Node struct {
	uuid          string
	nodeType      utils.Entity
	localAddress  string
	leader        LeaderNode
	nodes         map[string]string
	cancel        context.CancelFunc
	context       context.Context
	outputChannel chan string
}

func HandleNode(local string, remote string) {
	fmt.Println("Generating new node")
	// For cleanup
	context, cancel := context.WithCancel(context.Background())

	// Create the node
	node := Node{
		uuid:          uuid.New().String(),
		nodeType:      utils.Node,
		localAddress:  local,
		leader:        LeaderNode{uuid: "", leaderAddress: ""},
		nodes:         make(map[string]string),
		cancel:        cancel,
		context:       context,
		outputChannel: make(chan string, 100),
	}
	node.nodes[node.uuid] = node.localAddress

	// Listen for incoming
	go node.listen()

	// Set up Console
	var wait sync.WaitGroup
	wait.Add(1)
	go node.inputHandler(&wait)

	// Initialization
	if remote != "" {
		// Register with the leader
		node.outputChannel <- fmt.Sprintf("Attempting to register with Leader at: %s", remote)
		registerRequest, err := RegisterRequest(&node)
		if err != nil {
			fmt.Printf("Could not generate registration request: %v", err)
			node.Shutdown()
		}

		// Connect with leader
		conn, err := net.Dial(string(utils.TCP), remote)
		if err != nil {
			node.outputChannel <- fmt.Sprintf("Could not generate a connection with leader at %s on %s: %v", remote, node.localAddress, err)
			node.Shutdown()
		}
		// Send the request
		conn.Write(registerRequest)
		conn.Close()
	} else {
		node.nodeType = utils.Leader
		node.leader.uuid = node.uuid
		node.leader.leaderAddress = node.localAddress
	}

	wait.Wait()
	node.Shutdown()
}

func (node *Node) outputHandler() {
	for out := range node.outputChannel {
		select {
		case <-node.context.Done():
			return

		default:
			fmt.Println(out)
		}
	}
}

func (node *Node) inputHandler(wait *sync.WaitGroup) {
	defer wait.Done()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		select {
		case <-node.context.Done():
			return

		default:
			text := strings.ToLower(scanner.Text())
			node.outputChannel <- "You typed " + text
			if text == "exit" {
				node.context.Done()
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		node.outputChannel <- "Error reading from the console: " + err.Error()
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

	if node.nodeType != utils.Leader {
		conn, err := net.Dial(string(utils.TCP), node.leader.leaderAddress)
		if err != nil {
			node.outputChannel <- fmt.Sprintf("Could not contact leader %s :: %s", node.leader.uuid, node.leader.leaderAddress)
		}
		defer conn.Close()
	}
}

func (node *Node) printNodeList() {
	for id, a := range node.nodes {
		node.outputChannel <- fmt.Sprintf("    - %s :: %s", id, a)
	}
	node.outputChannel <- ""
}

func (node *Node) listen() {
	node.outputChannel <- "Initiating listener..."
	// Create the listener
	listener, err := net.Listen(string(utils.TCP), node.localAddress)
	if err != nil {
		node.outputChannel <- fmt.Sprintf("Failed to listen on %s: %v\n", node.localAddress, err)
		node.Shutdown()
	}
	defer listener.Close()

	node.outputChannel <- fmt.Sprintf("Listening on %s\n", node.localAddress)
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

func (node *Node) Shutdown() {
	// If current leader, ensure the other nodes are updated first
	if node.uuid == node.leader.uuid && len(node.nodes) > 0 {
		for k, v := range node.nodes {
			if k == node.uuid {
				continue
			}
			node.leader.uuid = k
			node.leader.leaderAddress = v
			node.nodeType = utils.Node
			break
		}
		delete(node.nodes, node.uuid)
		update, err := UpdateNodeListRequest(node)
		if err != nil {
			fmt.Println("No moar leadaaa")
		}
		node.sendMessage(update)
	}

	shutdown, err := ShutdownRequest(node)
	if err != nil {
		fmt.Println("Rip ig??")
	}
	node.sendMessage(shutdown)

	close(node.outputChannel)
	fmt.Println("Exiting node")
	node.cancel()
}
