package client

import (
	"fmt"
	"go_server/utils"
	"net"
	"os"
)

type Client struct {
	host string
	port uint16
}

func HandleClient(address string) {
	connection, err := net.Dial(string(utils.TCP), address)
	if err != nil {
		fmt.Printf("Failed to connect to %s: %v\n", address, err)
		os.Exit(1)
	}
	defer connection.Close()

	fmt.Printf("Connected to %s\n", address)
}
