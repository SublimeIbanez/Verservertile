package user

import (
	"bufio"
	"fmt"
	"go_server/server"
	"go_server/utils"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

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
					client.processServiceSelection(text)
				}

			case ServiceChosen:
				{

				}

			case InService:
				{
					if client.service.Address == "" {
						client.outputChannel <- "No address assigned to current service"
						continue
					}

					switch client.service.Service {
					case server.Chat:
						{
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		client.outputChannel <- "Error reading from the console: " + err.Error()
	}
}

func (client *Client) processServiceSelection(text string) {
	choice, err := strconv.Atoi(text)
	if err != nil {
		client.outputChannel <- "Invalid choice, you must make a selection [number]"
		return
	}
	if choice <= 0 || choice > len(client.services) {
		client.outputChannel <- "Invalid choice, please make a valid selection"
		return
	}
	i := 1
	for service := range client.services {
		if choice != i {
			i += 1
			continue
		}
		request, err := client.ServiceChoiceRequest(service)
		if err != nil {
			client.outputChannel <- "Could not generate service choice request: " + err.Error()
			return
		}

		conn, err := net.Dial(string(utils.TCP), client.remoteAddress)
		if err != nil {
			client.outputChannel <- "Could not send service choice request: " + err.Error()
			return
		}
		defer conn.Close()
		conn.Write(request)
		break
	}
	client.outputChannel <- "Awaiting response from the server..."
}
