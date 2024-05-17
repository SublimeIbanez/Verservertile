package user

import (
	"encoding/json"
	"errors"
	"go_server/common"
	"go_server/server"
	"go_server/utils"
	"slices"
)

func (client *Client) ServiceListRequest() ([]byte, error) {
	base := common.NewRequest(true, []string{""}, common.ServicesRequest)
	base.EntityType = utils.Client
	base.Uuid = client.uuid

	return base.Marshal(false)
}

func (client *Client) ServiceOperationRequest(ok bool, message []string, service string) ([]byte, error) {
	base := common.NewRequest(ok, message, common.ServiceOperation)
	base.EntityType = utils.Client
	base.Uuid = client.uuid

	return base.Marshal(true)
}

type ServiceChoice struct {
	Service server.Service
}

func (sc *ServiceChoice) Marshal() ([]byte, error) {
	if !slices.Contains(server.ServiceList, sc.Service) {
		return nil, errors.New("unsupported service provided")
	}

	return json.Marshal(sc)
}

func (client *Client) ServiceChoiceRequest(choice server.Service) ([]byte, error) {
	serviceChoice := ServiceChoice{
		Service: choice,
	}
	data, err := serviceChoice.Marshal()
	if err != nil {
		return nil, errors.New("Could not marshal new service choice request: " + err.Error())
	}

	base := common.NewRequest(true, []string{""}, common.ServiceChoice)
	base.EntityType = utils.Client
	base.Uuid = client.uuid
	base.Data = map[common.Directive]*json.RawMessage{
		common.ServiceChoice: (*json.RawMessage)(&data),
	}
	return base.Marshal(true)
}
