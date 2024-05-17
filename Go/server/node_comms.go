package server

import (
	"encoding/json"
	"fmt"
	"go_server/common"
	"go_server/utils"
	"io"
	"net"
	"slices"
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
	var baseMessage common.BaseMessage
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
			case common.Request:
				node.handleClientRequest(&baseMessage, &conn)
			case common.Response:
				node.handleClientResponse(baseMessage, conn)
			}
		}
	case utils.Node, utils.Leader:
		{
			switch baseMessage.Direction {
			case common.Request:
				node.handleNodeRequest(baseMessage, data)
			case common.Response:
				node.handleNodeResponse(baseMessage)
			}
		}
	default:
		node.outputChannel <- fmt.Sprintf("Invalid Source: %+v\n", baseMessage)
		return
	}
}

func (node *Node) handleClientRequest(base *common.BaseMessage, conn *net.Conn) {
	switch base.Directive {
	case common.ServicesRequest:
		{
			defer (*conn).Close()

			var message []string
			response, err := ServicesListResponse(node)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Could not generate service response: %v", err)
				message = append(message, fmt.Sprintf("Could not generate service response: %v", err))
				response, err = ErrorResponse(node, message, common.ServicesRequest)
				if err != nil {
					node.outputChannel <- fmt.Sprintf("Could not generate error response: %v", err)
					return
				}
			}
			(*conn).Write(response)
		}

	case common.ServiceOperation:
		{
			node.handleService(base, conn)
		}

	default:
		{
			defer (*conn).Close()
		}
	}
}

func (node *Node) handleClientResponse(base common.BaseMessage, conn net.Conn) {

}

func (node *Node) handleNodeRequest(base common.BaseMessage, data []byte) {
	switch base.Directive {
	case common.Register: // =========================================================================== Register
		{
			// Parse the request
			message := []string{""}
			ok := true
			var regRequest RegistrationRequest
			if _, k := node.nodes[base.Uuid]; k {
				// If the node already exists in the network
				ok = false
				message = append(message, "That node already exists on the network.")
			}

			if ok {
				// If the node doesn't exist
				if err := json.Unmarshal(*base.Data[common.Register], &regRequest); err != nil {
					ok = false
					message = append(message, "Could not unmarshal registration request data.\n")
				}
			}

			node.mtx.RLock()
			switch node.nodeType {
			case utils.Leader:
				{
					node.mtx.RUnlock()
					if ok {
						// Add new node to node list
						node.mtx.Lock()
						node.nodes[base.Uuid] = regRequest.Address
						for _, service := range regRequest.Services {
							// Make sure to not add a node that already exists within the list
							if nodeList, k := node.nodeServices[Service(service)]; k && slices.Contains(*nodeList, base.Uuid) {
								continue
							}
							// Add node's uuid to the services list
							if node.nodeServices[Service(service)] == nil {
								node.nodeServices[Service(service)] = &[]string{}
							}
							*node.nodeServices[Service(service)] = append(*node.nodeServices[Service(service)], base.Uuid)
						}

						node.mtx.Unlock()
						// Print results
						node.outputChannel <- fmt.Sprintf("New <node :: %s :: %s>\nNodeList:", base.Uuid, regRequest.Address)
						node.printNodeList()
						node.printNodeServiceList()
					}

					regResponse, err := RegisterResponse(node, ok, message)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not generate register response: %v", err)
						break
					}

					conn, err := net.Dial(string(utils.TCP), regRequest.Address)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not connect with node to register: %v", err)
						break
					}
					conn.Write(regResponse)
					conn.Close()

					update, err := UpdateNodeListRequest(node)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not generate update request: %v", err)
						break
					}

					node.sendMessage(update)
				}

			case utils.Node:
				{ // Only the leader should update, pass to leader
					node.mtx.RUnlock()
					conn, err := net.Dial(string(utils.TCP), node.leader.Address)
					if err != nil {
						node.outputChannel <- fmt.Sprintf("Could not connect with leader pass registration: %v", err)
						break
					}
					conn.Write(data)
					conn.Close()
				}
			}
		}

	case common.UpdateNodesList: // =========================================================================== UPDATE
		{
			// All but the leader handle this
			node.mtx.RLock()
			if node.nodeType == utils.Leader {
				node.mtx.RUnlock()
				return
			}
			node.mtx.RUnlock()

			node.outputChannel <- "Updating nodes list"

			var update UpdateRequest
			err := json.Unmarshal(*base.Data[common.UpdateNodesList], &update)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Could not unmarshal update request: %v", err)
				break
			}

			// Set the leader and update the Nodes list
			node.mtx.Lock()
			node.leader.Uuid = update.LeaderId
			node.leader.Address = update.LeaderAddress
			if node.leader.Uuid == node.uuid {
				node.nodeType = utils.Leader
				node.removeServiceNode(node.uuid)
				node.outputChannel <- "Updating node type to Leader"
			}
			node.nodes = update.Nodes
			node.nodeServices = update.Services
			node.mtx.Unlock()

			node.outputChannel <- "Nodes List updated:"
			node.printNodeList()
			node.outputChannel <- "Services List updated:"
			node.printNodeServiceList()
		}

	case common.Shutdown: // ================================================================================== Shutdown
		{
			// Only the leader manages this
			node.mtx.RLock()
			if node.nodeType != utils.Leader {
				node.mtx.RUnlock()
				return
			}

			// If the node isn't in the list, then simply return
			if _, ok := node.nodes[base.Uuid]; !ok {
				node.mtx.RUnlock()
				return
			}
			node.mtx.RUnlock()

			node.outputChannel <- fmt.Sprintf("Removing <%s::%s> from node lists", base.Uuid, node.nodes[base.Uuid])
			node.mtx.Lock()
			delete(node.nodes, base.Uuid)
			node.removeServiceNode(base.Uuid)
			node.mtx.Unlock()
			node.printNodeList()
			node.printNodeServiceList()

			update, err := UpdateNodeListRequest(node)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Unable to generate update request: %v", err)
				return
			}
			node.sendMessage(update)
		}
	}
}

func (node *Node) handleNodeResponse(base common.BaseMessage) {
	switch base.Directive {
	case common.Register:
		{
			// Get the information
			var leaderNode LeaderNode
			err := json.Unmarshal(*base.Data[common.Register], &leaderNode)
			if err != nil {
				node.outputChannel <- fmt.Sprintf("Could not unmarshal register response: %v", err)
				break
			}
			// Set the leader and update the Nodes list
			node.mtx.Lock()
			defer node.mtx.Unlock()
			node.leader.Uuid = leaderNode.Uuid
			node.leader.Address = leaderNode.Address
			if node.leader.Uuid == node.uuid {
				node.nodeType = utils.Leader
				node.removeServiceNode(node.uuid)
				node.outputChannel <- "Updating node type to Leader"
			}
		}
	}
}
