package protocol

import (
	"encoding/json"
	"errors"
	"go_server/utils"
)

type (
	// Direction (request/response)
	Direction uint8
	// Purpose of the message
	Directive uint8
)

const (
	// Direction - Nil is erroneous / incomplete base message
	NilDirection Direction = 0
	Request      Direction = 1
	Response     Direction = 2

	// Directive - Nil is erroneous / incomplete base message
	NilDirective    Directive = 0
	Register        Directive = 1
	Shutdown        Directive = 2
	UpdateNodesList Directive = 3
	ServiceRequest  Directive = 4
)

// `BaseMessage` is the only message being passed to and from entities within or connecting to the network
// Fields:
// - `Ok`: `boolean`, determines the status of the message
// - `Message`: `string` array for any message passed - usually accompanied by `!Ok`
// - `Direction`: `Request` or `Response`
// - `Directive`: Purpose of the message being passed
// - `EntityType`: Type of the entity passing the message - Usually `Client`, `Leader`, or `Node`
// - `Uuid`: Unique identifier of the entity passing the message
// - `Data`: Optional field which contains any extra information necessary for the purpose of the directive
type BaseMessage struct {
	Ok         bool
	Direction  Direction
	Directive  Directive
	EntityType utils.Entity
	Uuid       string
	Message    []string
	Data       map[Directive]*json.RawMessage
}

// Returns a partially implemented BaseMessage for a `Request`
func NewRequest(ok bool, message []string, directive Directive) BaseMessage {
	return BaseMessage{
		Ok:        ok,
		Message:   message,
		Direction: Request,
		Directive: directive,
	}
}

// Returns a partially implemented BaseMessage for a `Response`
func NewResponse(ok bool, message []string, directive Directive) BaseMessage {
	return BaseMessage{
		Ok:        ok,
		Message:   message,
		Direction: Response,
		Directive: directive,
	}
}

// Performs checks to ensure every relevant field is filled in the base message
// - `Ok` will either be true or false, there's no way to test that the field is properly applied
// - `Message` must always pass some message value, []string{""} is default
// - `Direction` must not be `NilDirection` - an erroneous detection method implemented
// - `Directive` must not be `NilDirective` - an erroneous detection value - ensures default value will be caught
// - `EntityType` must not be `NilEntity` - an erroneous detection value - ensures default value will be caught
// - `UUID` cannot be an empty string
// - `Data` is checked only if `dataAttached` was passed as `true`:
//   - If attaching data to the base message, ensure that the length of the data is not 0
func (base *BaseMessage) Marshal(dataAttached bool) ([]byte, error) {
	// Ok will either be true or false, there's no way to test that the field is properly applied
	// Must always pass some message value, []string{""} is default
	if base.Message == nil {
		return []byte{}, errors.New("message property in message passing cannot be nil")
	}
	// `NilDirection` is an erroneous detection value - ensures default value will be caught
	if base.Direction == NilDirection {
		return []byte{}, errors.New("invalid type, must be request or response")
	}
	// `NilDirective` is an erroneous detection value - ensures default value will be caught
	if base.Directive == NilDirective {
		return []byte{}, errors.New("invalid directive passed")
	}
	// `NilEntity` is an erroneous detection value - ensures default value will be caught
	if base.EntityType == utils.NilEntity {
		return []byte{}, errors.New("entity must pass a valid entity type")
	}
	// UUID cannot be an empty string
	if base.Uuid == "" {
		return []byte{}, errors.New("uuid cannot be blank")
	}
	// If attaching data to the base message, ensure that the length of the data is not 0
	if dataAttached && len(base.Data) == 0 {
		return []byte{}, errors.New("data must not be empty")
	}

	return json.Marshal(base)
}
