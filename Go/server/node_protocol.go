package server

import (
	"encoding/json"
	"go_server/protocol"
)

type UpdateRequest struct {
	LeaderId      string
	LeaderAddress string
	Nodes         map[string]string
}
type RegistrationRequest struct {
	Address string
}

func RegisterRequest(node *Node) ([]byte, error) {
	registrationRequest := RegistrationRequest{
		Address: node.localAddress,
	}
	register, err := json.Marshal(registrationRequest)
	if err != nil {
		return nil, err
	}
	baseMessage := protocol.BaseMessage{
		Ok:         true,
		Message:    []string{""},
		Type:       protocol.Request,
		Directive:  protocol.Register,
		SourceType: node.nodeType,
		Uuid:       node.uuid,
		Data: map[protocol.Directive]*json.RawMessage{
			protocol.Register: (*json.RawMessage)(&register),
		},
	}
	return json.Marshal(baseMessage)
}

func RegisterResponse(node *Node, ok bool, message []string) ([]byte, error) {
	registrationResponse := UpdateRequest{
		LeaderId:      node.leader.uuid,
		LeaderAddress: node.leader.leaderAddress,
		Nodes:         node.nodes,
	}
	register, err := json.Marshal(registrationResponse)
	if err != nil {
		return nil, err
	}
	baseMessage := protocol.BaseMessage{
		Ok:         ok,
		Message:    message,
		Type:       protocol.Response,
		Directive:  protocol.Register,
		SourceType: node.nodeType,
		Uuid:       node.uuid,
		Data: map[protocol.Directive]*json.RawMessage{
			protocol.Register: (*json.RawMessage)(&register),
		},
	}
	return json.Marshal(baseMessage)
}

func ShutdownRequest(node *Node) ([]byte, error) {
	base := protocol.BaseMessage{
		Ok:         true,
		Message:    []string{""},
		Type:       protocol.Request,
		Directive:  protocol.Shutdown,
		SourceType: node.nodeType,
		Uuid:       node.uuid,
	}

	return json.Marshal(base)
}

func UpdateNodeListRequest(node *Node) ([]byte, error) {
	updateResponse := UpdateRequest{
		LeaderId:      node.leader.uuid,
		LeaderAddress: node.leader.leaderAddress,
		Nodes:         node.nodes,
	}
	update, err := json.Marshal(updateResponse)
	if err != nil {
		return nil, err
	}

	base := protocol.BaseMessage{
		Ok:         true,
		Message:    []string{""},
		Type:       protocol.Request,
		Directive:  protocol.UpdateNodesList,
		SourceType: node.nodeType,
		Uuid:       node.uuid,
		Data: map[protocol.Directive]*json.RawMessage{
			protocol.UpdateNodesList: (*json.RawMessage)(&update),
		},
	}

	return json.Marshal(base)
}
