package user

import (
	"context"
	"fmt"
	"go_server/server"
	"go_server/utils"
	"net"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type State uint8

const (
	AwaitingServiceChoice State = 0
	ServiceChosen         State = 1
	InService             State = 2
)

type SelectedService struct {
	Service server.Service
	Address string
}

type Client struct {
	uuid          string
	localAddress  string
	remoteAddress string
	outputChannel chan string
	inputChannel  chan string
	services      map[server.Service]*[]string // service -> addresses
	service       SelectedService
	state         State
	mtx           sync.RWMutex
	cancel        context.CancelFunc
	context       context.Context
}

func HandleClient(local string, remote string) {
	fmt.Println("Generating new Client")
	// For cleanup
	context, cancel := context.WithCancel(context.Background())

	client := Client{
		uuid:          strings.ReplaceAll(uuid.NewString(), "-", ""),
		localAddress:  local,
		remoteAddress: remote,
		outputChannel: make(chan string, 100),
		inputChannel:  make(chan string, 100),
		services:      make(map[server.Service]*[]string),
		cancel:        cancel,
		context:       context,
	}
	go client.listen()

	var wait sync.WaitGroup
	wait.Add(1)
	go client.inputHandler(&wait)

	// Make a request to get all of the offered services
	services, err := client.ServiceListRequest()
	if err != nil {
		fmt.Printf("Could not make a request for services")
		client.Shutdown()
	}
	conn, err := net.Dial(string(utils.TCP), client.remoteAddress)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", client.remoteAddress, err)
		client.Shutdown()
	}
	conn.Write(services)
	client.handleConnection(conn)

	wait.Wait()
	client.Shutdown()
}

func (client *Client) listen() {
	client.outputChannel <- "Initiating listener..."
	// Create the listener
	listener, err := net.Listen(string(utils.TCP), client.localAddress)
	if err != nil {
		client.outputChannel <- fmt.Sprintf("Failed to listen on %s: %v\n", client.localAddress, err)
		client.Shutdown()
	}
	defer listener.Close()

	client.outputChannel <- fmt.Sprintf("Listening on %s\n", client.localAddress)
	go client.outputHandler()

	// Listen for incoming connections
	for {
		select {
		case <-client.context.Done():
			client.outputChannel <- "Stopping listener"
			return

		default:
			connection, err := listener.Accept()
			if err != nil {
				client.outputChannel <- fmt.Sprintf("Failed to accept connection: %v\n", err)
				continue
			}
			go client.handleConnection(connection)
		}
	}
}

func (client *Client) printServiceList() {
	client.outputChannel <- "Service List:"
	i := 1
	for service := range client.services {
		client.outputChannel <- fmt.Sprintf("  %d. %s", i, service)
		i += 1
	}
	client.state = AwaitingServiceChoice
	client.outputChannel <- "-- Please make a selection -- "
	client.outputChannel <- ""
}

func (client *Client) Shutdown() {

}
