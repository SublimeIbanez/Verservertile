package user

import (
	"bufio"
	"context"
	"fmt"
	"go_server/server"
	"go_server/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type State uint8

const (
	AwaitingServiceChoice State = 0
	ServiceChosen         State = 1
)

type Client struct {
	uuid          string
	localAddress  string
	remoteAddress string
	outputChannel chan string
	services      map[server.Service]*[]string // service -> addresses
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
		cancel:        cancel,
		context:       context,
		services:      make(map[server.Service]*[]string),
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

func (client *Client) outputHandler() {
	for out := range client.outputChannel {
		select {
		case <-client.context.Done():
			return

		default:
			fmt.Println(out)
		}
	}

}

func (client *Client) inputHandler(wait *sync.WaitGroup) {
	defer wait.Done()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		select {
		case <-client.context.Done():
			return

		default:
			text := strings.ToLower(scanner.Text())
			if text == "/exit" {
				client.context.Done()
				return
			}
			switch client.state {
			case AwaitingServiceChoice:
				{
					choice, err := strconv.Atoi(text)
					if err != nil {
						client.outputChannel <- "Invalid choice, you must make a selection [number]"
						continue
					}
					if choice <= 0 || choice > len(client.services) {
						client.outputChannel <- "Invalid choice, please make a valid selection"
						continue
					}
					i := 1
					for service := range client.services {
						if choice != i {
							continue
						}
						request, err := client.ServiceChoiceRequest(service)
						if err != nil {
							client.outputChannel <- "Could not generate service choice request: " + err.Error()
						}

						conn, err := net.Dial(string(utils.TCP), client.remoteAddress)
						if err != nil {
							client.outputChannel <- "Could not send service choice request: " + err.Error()
						}
						defer conn.Close()
						conn.Write(request)
						break
					}
					client.outputChannel <- "Awaiting response from the server..."
				}

			}
		}
	}

	if err := scanner.Err(); err != nil {
		client.outputChannel <- "Error reading from the console: " + err.Error()
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
