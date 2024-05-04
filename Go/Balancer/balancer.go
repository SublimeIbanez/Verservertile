package balancer

import (
	"fmt"
	"go_server/node"
	"go_server/utils"
	"net"
	"os"
	"strconv"
	"strings"
)

type BalancerNode struct {
	host  string
	port  uint16
	nodes []node.Node
}

func (balancer BalancerNode) newBalancer(address string) BalancerNode {
	var prt, err = strconv.Atoi(strings.Split(address, ":")[1])
	if err != nil {
		fmt.Printf("How did you dun get past the initial parse??!!?! %s\n", err)
	}

	return BalancerNode{
		host: strings.Split(address, ":")[0],
		port: uint16(prt),
	}
}

func HandleBalancer(address string) {
	listener, err := net.Listen(string(utils.TCP), address)
	if err != nil {
		fmt.Printf("Failed to listen on %s: %v", address, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("Listening on %s\n", address)
	for {
		_, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
	}
}
