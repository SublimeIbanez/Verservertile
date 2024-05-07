package client

import (
	"context"
	"fmt"
	"go_server/utils"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Client struct {
	uuid          string
	localAddress  string
	remoteAddress string
	outputChannel chan string
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
		cancel:        cancel,
		context:       context,
	}
	client.listen()

	conn, err := net.Dial(string(utils.TCP), client.remoteAddress)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", client.remoteAddress, err)
		os.Exit(1)
	}
	defer conn.Close()

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

func (client *Client) handleConnection(conn net.Conn) {
	defer conn.Close()

}

func (client *Client) outputHandler() {

}

func (client *Client) inputHandler() {

}

func (client *Client) Shutdown() {

}
