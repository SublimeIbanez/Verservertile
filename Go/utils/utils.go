package utils

type Mode uint8
type Protocol string

const (
	Server   Mode = 0
	Client   Mode = 1
	Balancer Mode = 2
)

var ValidModes = map[string]Mode{
	"server":   Server,
	"client":   Client,
	"balancer": Balancer,
}

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)
