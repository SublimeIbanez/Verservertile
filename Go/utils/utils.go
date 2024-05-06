package utils

type Entity uint8
type ConnType string

const (
	Leader Entity = 0
	Node   Entity = 1
	Client Entity = 2
)

var EntityTypes = map[string]Entity{
	"leader": Leader,
	"node":   Node,
	"client": Client,
}

const (
	TCP ConnType = "tcp"
	UDP ConnType = "udp"
)

const BUFFER_SIZE = 1024
