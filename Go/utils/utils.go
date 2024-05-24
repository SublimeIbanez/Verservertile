package utils

type Entity uint8
type ConnType string

const (
	// NilEntity ensures the default value is never passed
	NilEntity Entity = 0
	Leader    Entity = 1
	Node      Entity = 2
	Candidate Entity = 3
	Client    Entity = 4

	// NilConnType ensures the default value is never passed
	NilConnType ConnType = ""
	TCP         ConnType = "tcp"
	UDP         ConnType = "udp"
)

// Translates string to entity types
var EntityTypes = map[string]Entity{
	"leader": Leader,
	"node":   Node,
	"client": Client,
}

const (
	// Maximum buffer size in bytes for message passing
	BUFFER_SIZE = 1_024
)
