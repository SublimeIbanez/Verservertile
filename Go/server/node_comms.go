package server

import (
	"encoding/json"
	"fmt"
	"go_server/protocol"
	"go_server/utils"
	"io"
	"net"
)

func (node *Node) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, utils.BUFFER_SIZE)
	size, err := conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			node.outputChannel <- fmt.Sprintf("Error reading from connection: %v\n", err)
		}
		return
	}

	data := buffer[:size]
	var baseMessage protocol.BaseMessage
	err = json.Unmarshal(data, &baseMessage)
	if err != nil {
		node.outputChannel <- fmt.Sprintf("Error attempting to parse JSON: %v\n", err)
		return
	}

	if !baseMessage.Ok {
		node.outputChannel <- fmt.Sprintf("Message has returned an error:\n  %v", baseMessage.Message)
		return
	}

	switch baseMessage.EntityType {
	case utils.Client:
		{
			switch baseMessage.Direction {
			case protocol.Request:
				node.handleClientRequest(baseMessage)
			case protocol.Response:
				node.handleClientResponse(baseMessage)
			}
		}
	case utils.Node, utils.Leader:
		{
			switch baseMessage.Direction {
			case protocol.Request:
				node.handleNodeRequest(baseMessage, data)
			case protocol.Response:
				node.handleNodeResponse(baseMessage)
			}
		}
	default:
		node.outputChannel <- fmt.Sprintf("Invalid Source: %+v\n", baseMessage)
		return
	}
}

func (node *Node) handleClientRequest(base protocol.BaseMessage) {

}

func (node *Node) handleClientResponse(base protocol.BaseMessage) {

}

func (node *Node) handleNodeRequest(base protocol.BaseMessage, data []byte) {
	switch base.Directive {
	case protocol.Register: // ======= Register
		{
			// Parse the request
			message := []string{""}
			ok := true
			var address string
			if _, k := node.nodes[base.Uuid]; k {
				// If the node already exists in the network
				ok = false
				message = append(message, "That node already exists on the network.")
			}

			if ok {
				// If the node doesn't exist
				var register RegistrationRequest
				if err := json.Unmarshal(*base.Data[protocol.Register], &register); err != nil {
					ok = false
					message = append(message, "Could not unmarshal registration request data.\n")
				} else {
					address = register.Address
				}
			}

			switch node.nodeType {
			case utils.Leader:
				{
					if ok {
						node.nodes[base.Uuid] = address
						node.outputChannel <- fmt.Sprintf("New <node :: %s :: %s>\nNodeList:", base.Uuid, address)
						node.printNodeList()
					}

					register, err := RegisterResponse(node, ok, message)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not generate register response: %v", err)
						break
					}

					conn, err := net.Dial(string(utils.TCP), address)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not connect with node to register: %v", err)
						break
					}
					conn.Write(register)
					conn.Close()

					update, err := UpdateNodeListRequest(node)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not generate update request: %v", err)
						break
					}

					go node.sendMessage(update)

				}

			case utils.Node:
				{ // Only the leader should update, pass to leader
					conn, err := net.Dial(string(utils.TCP), node.leader.address)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not connect with leader pass registration: %v", err)
						break
					}
					conn.Write(data)
					conn.Close()
				}
			}
		}

	case protocol.UpdateNodesList: // ======= UPDATE
		{
			// All but the leader handle this
			if node.nodeType == utils.Leader {
				return
			}

			node.outputChannel <- "Updating nodes list"

			var update UpdateRequest
			err := json.Unmarshal(*base.Data[protocol.UpdateNodesList], &update)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Could not unmarshal update request: %v", err)
				break
			}

			// Set the leader and update the Nodes list
			node.leader.uuid = update.LeaderId
			node.leader.address = update.LeaderAddress
			if node.leader.uuid == node.uuid {
				node.nodeType = utils.Leader
				node.outputChannel <- "Updating node type to Leader"
			}

			node.nodes = update.Nodes
			node.outputChannel <- "Nodes List updated:"
			node.printNodeList()
		}

	case protocol.Shutdown: // ======= Shutdown
		{
			// Only the leader manages this
			if node.nodeType != utils.Leader {
				return
			}

			// If the node isn't in the list, then simply return
			if _, ok := node.nodes[base.Uuid]; !ok {
				return
			}

			node.outputChannel <- fmt.Sprintf("Removing <%s :: %s> from node list", base.Uuid, node.nodes[base.Uuid])
			delete(node.nodes, base.Uuid)

			update, err := UpdateNodeListRequest(node)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Unable to generate udpate request: %v", err)
				return
			}
			node.sendMessage(update)
		}
	}
}

func (node *Node) handleNodeResponse(base protocol.BaseMessage) {
	switch base.Directive {
	case protocol.Register:
		{
			// Get the information
			var register UpdateRequest
			err := json.Unmarshal(*base.Data[protocol.Register], &register)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Could not unmarshal register response: %v", err)
				break
			}
			// Set the leader and update the Nodes list
			node.leader.uuid = register.LeaderId
			node.leader.address = register.LeaderAddress
			if node.leader.uuid == node.uuid {
				node.nodeType = utils.Leader
				node.outputChannel <- "Updating node type to Leader"
			}

			node.nodes = register.Nodes
			node.outputChannel <- "Nodes List updated:"
			node.printNodeList()
		}
	}
}
