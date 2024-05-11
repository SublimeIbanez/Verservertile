package server

import (
	"encoding/json"
	"errors"
	"go_server/common"
)

type UpdateRequest struct {
	LeaderId      string
	LeaderAddress string
	Nodes         map[string]string
	Services      map[string]*[]string
}

func (ur *UpdateRequest) Marshal() ([]byte, error) {
	if ur.LeaderId == "" {
		return nil, errors.New("leader id cannot be empty")
	}
	if ur.LeaderAddress == "" {
		return nil, errors.New("leader address cannot be empty")
	}
	if ur.Nodes == nil {
		return nil, errors.New("nodes cannot be a nil value")
	}
	if ur.Services == nil {
		return nil, errors.New("services cannot be a nil value")
	}

	return json.Marshal(ur)
}

type RegistrationRequest struct {
	Address  string
	Services []string
}

func (rr *RegistrationRequest) Marshal() ([]byte, error) {
	if rr.Address == "" {
		return nil, errors.New("invalid address, node must have an address")
	}
	if len(rr.Services) == 0 {
		return nil, errors.New("services cannot be empty")
	}
	return json.Marshal(rr)
}

func (leader *LeaderNode) Marshal() ([]byte, error) {
	if leader.Uuid == "" {
		return nil, errors.New("uuid cannot be empty")
	}
	if leader.Address == "" {
		return nil, errors.New("invalid address, leader must have an address")
	}

	return json.Marshal(leader)
}

func ErrorResponse(node *Node, message []string, directive common.Directive) ([]byte, error) {
	base := common.NewResponse(false, message, directive)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid

	return base.Marshal(false)
}

func RegisterRequest(node *Node) ([]byte, error) {
	regReq := RegistrationRequest{
		Address:  node.address,
		Services: node.serviceList,
	}
	register, err := regReq.Marshal()
	if err != nil {
		return nil, err
	}

	base := common.NewRequest(true, []string{""}, common.Register)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[common.Directive]*json.RawMessage{
		common.Register: (*json.RawMessage)(&register),
	}

	return base.Marshal(true)
}

func RegisterResponse(node *Node, ok bool, message []string) ([]byte, error) {
	// Complete the struct
	leaderNode := LeaderNode{
		Uuid:    node.leader.Uuid,
		Address: node.leader.Address,
	}
	// Marshal - ensure check for errors
	leader, err := leaderNode.Marshal()
	if err != nil {
		return nil, err
	}
	// Generate a new baseMessage and fill out to the values
	base := common.NewResponse(ok, message, common.Register)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[common.Directive]*json.RawMessage{
		common.Register: (*json.RawMessage)(&leader),
	}

	// Return marshal, checking for errors
	return base.Marshal(true)
}

func ShutdownRequest(node *Node) ([]byte, error) {
	base := common.NewRequest(true, []string{""}, common.Shutdown)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid

	return base.Marshal(false)
}

func UpdateNodeListRequest(node *Node) ([]byte, error) {
	updateReq := UpdateRequest{
		LeaderId:      node.leader.Uuid,
		LeaderAddress: node.leader.Address,
		Nodes:         node.nodes,
		Services:      node.nodeServices,
	}

	update, err := updateReq.Marshal()
	if err != nil {
		return nil, err
	}

	base := common.NewRequest(true, []string{""}, common.UpdateNodesList)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[common.Directive]*json.RawMessage{
		common.UpdateNodesList: (*json.RawMessage)(&update),
	}

	return base.Marshal(true)
}

func ServicesResponse(node *Node) ([]byte, error) {
	serviceResponse := common.ServiceResponse{
		ServicesList: node.nodeServices,
	}

	response, err := serviceResponse.Marshal()
	if err != nil {
		return nil, err
	}

	base := common.NewResponse(true, []string{""}, common.ServiceRequest)
	base.EntityType = node.nodeType
	base.Uuid = node.uuid
	base.Data = map[common.Directive]*json.RawMessage{
		common.ServiceRequest: (*json.RawMessage)(&response),
	}

	return base.Marshal(true)
}
