package server

import (
	"fmt"
	"go_server/utils"
	"net"
	"os"

	"github.com/google/uuid"
)

var balancer Balancer

type Balancer struct {
	uuid    string
	address string
	nodes   map[string]string
}

func HandleBalancer(port uint16) {
	// Get the address
	ipAddress, err := utils.GetIp()
	if err != nil {
		fmt.Printf("Could not get local address: %s\n", err)
		os.Exit(1)
	}
	address := fmt.Sprintf("%s:%d", ipAddress, port)

	balancer = Balancer{
		uuid:    uuid.New().String(),
		address: address,
	}

	listener, err := net.Listen(string(utils.TCP), address)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v", address, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go balancer.handleConnection(connection)
	}
}

func (balancer *Balancer) handleConnection(connection net.Conn) {
	defer connection.Close()

	buffer := make([]byte, utils.BUFFER_SIZE)
	length, err := connection.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading from connection: %v\n", err)
		return
	}

	message := string(buffer[:length])
	fmt.Println(message)
}
