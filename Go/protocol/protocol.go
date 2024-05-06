package protocol

import (
	"encoding/json"
	"go_server/utils"
)

type (
	Prop      string
	Direction string
	Directive string
)

const (
	Type    Prop = "type"
	Source  Prop = "source"
	Uuid    Prop = "uuid"
	Address Prop = "address"
	Nodes   Prop = "nodes"
	Leader  Prop = "leader"
	Ok      Prop = "ok"
	Message Prop = "message"

	Request  Direction = "request"
	Response Direction = "response"

	Register        Directive = "register"
	Shutdown        Directive = "shutdown"
	UpdateNodesList Directive = "updatenodeslist"
)

type BaseMessage struct {
	Ok         bool
	Message    []string
	Type       Direction
	Directive  Directive
	SourceType utils.Entity
	Uuid       string
	Data       map[Directive]*json.RawMessage
}
