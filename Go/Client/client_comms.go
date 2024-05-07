package client

import (
	"encoding/json"
	"go_server/protocol"
	"go_server/utils"
)

// Ok         bool
// Message    []string
// Type       Direction
// Directive  Directive
// SourceType utils.Entity
// Uuid       string
// Data       map[Directive]*json.RawMessage
func (client *Client) ServiceRequest() ([]byte, error) {
	base := protocol.BaseMessage{
		Ok:         true,
		Message:    []string{""},
		Direction:  protocol.Request,
		Directive:  protocol.ServiceRequest,
		EntityType: utils.Client,
		Uuid:       client.uuid,
		Data:       nil,
	}

	return json.Marshal(base)
}
