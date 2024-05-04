package server

import "encoding/json"

type Property string
type Directive string

const (
	Type    Property = `json:"type"`
	Uuid    Property = `json:"uuid"`
	Address Property = `json:"address"`
	Nodes   Property = `json:"nodes"`
	Leader  Property = `json:"leader"`
	Ok      Property = `json:"ok"`
	Message Property = `json:"message"`

	Request  Directive = "request"
	Response Directive = "response"
	Register Directive = "register"
)

var PropertyMap = map[string]Property{
	"type":    Type,
	"uuid":    Uuid,
	"address": Address,
	"nodes":   Nodes,
	"leader":  Leader,
	"ok":      Ok,
	"message": Message,
}
var DirectiveMap = map[string]Directive{
	"register": Register,
}

type BaseMessage struct {
	Fields map[Property]interface{}
}

func (baseMessage BaseMessage) Marshal() ([]byte, error) {
	fieldMap := make(map[string]interface{})
	for key, value := range baseMessage.Fields {
		fieldMap[string(key)] = value
	}
	return json.Marshal(fieldMap)
}

func (baseMessage BaseMessage) UnMarshal(msg []byte) error {
	fieldMap := make(map[string]interface{})
	if err := json.Unmarshal(msg, &fieldMap); err != nil {
		return err
	}

	baseMessage.Fields = make(map[Property]interface{})
	for key, value := range fieldMap {
		baseMessage.Fields[Property(key)] = value
	}

	return nil
}

func (node *Node) RegisterRequest() ([]byte, error) {
	msg := BaseMessage{
		Fields: map[Property]interface{}{
			Type:    string(Register),
			Uuid:    node.uuid,
			Address: node.address,
		},
	}
	return msg.Marshal()
}

func (balancer *Balancer) RegisterResponse(message string) ([]byte, error) {
	baseMessage := BaseMessage{
		Fields: map[Property]interface{}{
			Type:    string(Register),
			Uuid:    balancer.uuid,
			Nodes:   balancer.nodes,
			Ok:      true,
			Message: message,
		},
	}
	return baseMessage.Marshal()
}
