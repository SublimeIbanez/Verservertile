package server

import (
	"encoding/json"
	"errors"
	"go_server/protocol"
)

type UpdateRequest struct {
	LeaderId      string
	LeaderAddress string
	Nodes         map[string]string
}

func (ur *UpdateRequest) Marshal() ([]byte, error) {
	if ur.LeaderId == "" {
		return []byte{}, errors.New("leader id cannot be empty")
	}
	if ur.LeaderAddress == "" {
		return []byte{}, errors.New("leader address cannot be empty")
	}
	if ur.Nodes == nil {
		return []byte{}, errors.New("nodes cannot be a nil value")
	}

	return json.Marshal(ur)
}

type RegistrationRequest struct {
	Address string
}

func (rr *RegistrationRequest) Marshal() ([]byte, error) {
	if rr.Address == "" {
		return []byte{}, errors.New("invalid address, node must have an address")
	}
	return json.Marshal(rr)
}

func RegisterRequest(node *Node) ([]byte, error) {
	regReq := RegistrationRequest{
		Address: node.address,
	}
	register, err := regReq.Marshal()
	if err != nil {
		return nil, err
	}

	base := protocol.NewRequest(true, []string{""}, protocol.Register)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[protocol.Directive]*json.RawMessage{
		protocol.Register: (*json.RawMessage)(&register),
	}

	return base.Marshal(true)
}

func RegisterResponse(node *Node, ok bool, message []string) ([]byte, error) {
	// Complete the struct
	regRes := UpdateRequest{
		LeaderId:      node.leader.uuid,
		LeaderAddress: node.leader.address,
		Nodes:         node.nodes,
	}
	// Marshal - ensure check for errors
	register, err := regRes.Marshal()
	if err != nil {
		return nil, err
	}
	// Generate a new baseMessage and fill out to the values
	base := protocol.NewResponse(ok, message, protocol.Register)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[protocol.Directive]*json.RawMessage{
		protocol.Register: (*json.RawMessage)(&register),
	}

	// Return marshal, checking for errors
	return base.Marshal(true)
}

func ShutdownRequest(node *Node) ([]byte, error) {
	base := protocol.NewRequest(true, []string{""}, protocol.Shutdown)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid

	return base.Marshal(false)
}

func UpdateNodeListRequest(node *Node) ([]byte, error) {
	updateReq := UpdateRequest{
		LeaderId:      node.leader.uuid,
		LeaderAddress: node.leader.address,
		Nodes:         node.nodes,
	}

	update, err := updateReq.Marshal()
	if err != nil {
		return nil, err
	}

	base := protocol.NewRequest(true, []string{""}, protocol.UpdateNodesList)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[protocol.Directive]*json.RawMessage{
		protocol.UpdateNodesList: (*json.RawMessage)(&update),
	}

	return base.Marshal(true)
}
