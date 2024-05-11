package user

import (
	"go_server/common"
	"go_server/utils"
)

func (client *Client) ServiceRequest() ([]byte, error) {
	base := common.NewRequest(true, []string{""}, common.ServiceRequest)
	base.EntityType = utils.Client
	base.Uuid = client.uuid

	return base.Marshal(false)
}

// Service -> Maps the level to the service
type ServiceOperation struct {
	Levels   uint8
	Services map[uint8][]Instruction
}

// ==Chat
// 1. Join Channel
// 2. Join Private
// 3. Create Channel
// 4. Create Private
// Back

// ==Join Channel
// Input channel name:

// Thing3:
// Put this input:
// --> user puts input
// Put this other input:
// --> user puts in other input
// <-- Returns back to Menu2 or performs next action

// Instruction: 0000 0000 0000 0000 0000 0000 0000 0000
// - 0bx0000. == Back one level
// - 0bx0001. == Display content -- Back Command: Display[0]
// - 0bx0011. == Await user selection --> Perform next instruction -- Back Command: Display[0]
// - 0bx0011. == Display content --> Take user input (Display: []string -> foreach user inputs string) --> Perform Next Instruction
// - 0bx0100. == Display content --> Take user input (Display: []string -> foreach user inputs string) --> Perform Previous Instruction
// - 0bx0101. == Send/Receive :: Argument Prefix: Display[0] -> Location: Display[1] -> Update Interval: Display[2]
// - 0bx0111. == Send/Recieve :: Update -- Display[0]
// - 0bx1000. == Load from FS
// - 0bx1001. == Save to FS
// - 0b.1111x
type Instruction struct {
	ServiceId   string
	Level       uint8
	Title       string
	Instruction uint32
	Commands    map[string]Instruction
	Display     []string
	Input       []string
}

func (client *Client) ServiceOperationRequest(ok bool, message []string, service string) ([]byte, error) {
	base := common.NewRequest(ok, message, common.ServiceOperation)
	base.EntityType = utils.Client
	base.Uuid = client.uuid

	return base.Marshal(true)
}
