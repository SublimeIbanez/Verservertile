package utils

import "net"

type Mode uint8
type Protocol string

const (
	Balancer Mode = 0
	Server   Mode = 1
	Client   Mode = 2
)

var ValidModes = map[string]Mode{
	"balancer": Balancer,
	"server":   Server,
	"client":   Client,
}

const (
	TCP Protocol = "tcp"
	UDP Protocol = "udp"
)

const BUFFER_SIZE = 1024

func GetIp() (string, error) {
	// Get local IP address
	conn, err := net.Dial(string(UDP), "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr).String()
	return ipAddress, nil
}
