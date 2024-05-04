package server

import (
	"context"
	"fmt"
	"go_server/utils"
	"net"
	"os"

	"github.com/google/uuid"
)

var node Node

type Node struct {
	uuid     string
	address  string
	balancer Balancer
	leader   bool
	nodes    map[string]string
	cancel   context.CancelFunc
}

func HandleNode(balancerAddress string, port uint16) {
	// Get the node's IP address
	ipAddress, err := utils.GetIp()
	if err != nil {
		fmt.Printf("Could not get local address: %s\n", err)
		os.Exit(1)
	}
	nodeAddress := fmt.Sprintf("%s:%d", ipAddress, port)

	// For cleanup
	context, cancel := context.WithCancel(context.Background())

	defer cancel()
	// Create this node
	node = Node{
		uuid:     uuid.New().String(),
		address:  nodeAddress,
		balancer: Balancer{address: balancerAddress, nodes: nil},
		cancel:   cancel,
	}

	// Listen for incoming
	go node.listen(context)

	// Register with the balancer
	fmt.Printf("Attempting to register with Balancer at: %s", node.balancer.address)
	registration, err := node.RegisterRequest()
	if err != nil {
		fmt.Printf("Could not generate registration: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.Dial(string(utils.TCP), balancerAddress)
	if err != nil {
		fmt.Printf("Could not generate a connection with balancer at %s on %s: %v\n", node.balancer.address, node.address, err)
		os.Exit(1)
	}
	defer conn.Close()
	conn.Write(registration)
}

func (node *Node) listen(context context.Context) {
	// Create the listener
	listener, err := net.Listen(string(utils.TCP), node.address)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v", node.address, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", node.address)

	// Listen for incoming connections
	for {
		select {
		case <-context.Done():
			fmt.Println("Stopping listener")
			return

		default:
			connection, err := listener.Accept()
			if err != nil {
				fmt.Printf("Failed to accept connection: %v\n", err)
				continue
			}
			go handleNodeConnection(connection)
		}
	}
}

func handleNodeConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, utils.BUFFER_SIZE)
	size, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		node.cancel()
		return
	}

	data := buffer[:size]
	var baseMessage BaseMessage
	err = baseMessage.UnMarshal(data)
	if err != nil {
		fmt.Printf("Error attempting to parse JSON: %v\n", err)
		node.cancel()
		return
	}

	if ok, valid := baseMessage.Fields[Ok].(bool); valid {
		if !ok {
			fmt.Printf("Message has returned an error:\n  %s", baseMessage.Fields[Message])
			node.cancel()
			return
		}
	} else {
		fmt.Printf("Could not parse OK status")
		node.cancel()
		return
	}

	t := baseMessage.Fields[Type]

	switch t {
	case Register:
		{
			var missingFields []string
			uuid, uuidOk := baseMessage.Fields[Uuid].(string)
			leader, leaderOk := baseMessage.Fields[Leader].(bool)
			nodes, nodesOk := baseMessage.Fields[Nodes].(map[string]string)

			if !uuidOk {
				missingFields = append(missingFields, string(Uuid))
			}
			if !leaderOk {
				missingFields = append(missingFields, string(Leader))
			}
			if !nodesOk {
				missingFields = append(missingFields, string(Nodes))
			}

			if len(missingFields) > 0 {
				fmt.Printf("Could not correctly parse message, missing fields: %v\n", missingFields)
				node.cancel()
				return
			}

			node.balancer.uuid = uuid
			node.leader = leader
			node.nodes = nodes
		}

	default:
		fmt.Printf("wut")
	}

}
