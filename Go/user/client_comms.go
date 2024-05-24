package user

import (
	"encoding/json"
	"fmt"
	"go_server/common"
	"go_server/server"
	"go_server/utils"
	"io"
	"net"
)

func (client *Client) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, utils.BUFFER_SIZE)
	size, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			client.outputChannel <- fmt.Sprintf("Error reading from connection: %v\n", err)
		}
		return
	}

	data := buffer[:size]
	var baseMessage common.BaseMessage
	err = json.Unmarshal(data, &baseMessage)
	if err != nil {
		client.outputChannel <- fmt.Sprintf("Error attempting to parse JSON: %v\n", err)
		return
	}

	if !baseMessage.Ok {
		client.outputChannel <- fmt.Sprintf("Message has returned an error:\n  %v", baseMessage.Message)
		return
	}

	switch baseMessage.EntityType {
	case utils.Node, utils.Candidate, utils.Leader:
		{
			switch baseMessage.Direction {
			case common.Request:
				client.handleRequest(baseMessage)
			case common.Response:
				client.handleResponse(baseMessage)
			}
		}
	default:
		client.outputChannel <- fmt.Sprintf("Invalid Source: %+v\n", baseMessage)
		return
	}

}

func (client *Client) handleRequest(base common.BaseMessage) {

}

func (client *Client) handleResponse(base common.BaseMessage) {
	switch base.Directive {
	case common.ServicesRequest:
		{
			var servicesResponse server.ServiceListResponse
			err := json.Unmarshal(*base.Data[common.ServicesRequest], &servicesResponse)
			if err != nil {
				client.outputChannel <- fmt.Sprintf("Could not unmarshal the service response: %v", err)
			}

			client.mtx.Lock()
			client.services = servicesResponse.ServicesList
			client.mtx.Unlock()
			client.outputChannel <- "Updated service list..."
			client.printServiceList()
		}

	default:
	}

}
