package balancer

import (
	"fmt"
	"go_server/node"
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
	fmt.Println("baransu")
}
