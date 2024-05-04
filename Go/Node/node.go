package node

import (
	"fmt"
	"go_server/client"
	"go_server/utils"
	"net"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	host       string
	port       uint16
	leader     bool
	nodeList   []Node
	clientList []client.Client
}

func (node Node) newNode(address string) Node {
	var prt, err = strconv.Atoi(strings.Split(address, ":")[1])
	if err != nil {
		fmt.Printf("How did you dun get past the initial parse??!!?! %s\n", err)
	}

	return Node{
		host: strings.Split(address, ":")[0],
		port: uint16(prt),
	}
}

func HandleNode(address string) {
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
		go handleNodeConnection(connection)
	}
}

func handleNodeConnection(connection net.Conn) {

}
