package server

import (
	"fmt"
	"slices"
)

var Services = []string{
	"chat",
}

func (node *Node) initServices() {
	node.updateServices(Services, node.uuid)
}

func (node *Node) updateServices(services []string, uuid string) {
	for _, service := range services {
		if node.nodeServices[service] == nil {
			node.nodeServices[service] = &[]string{}
		}
		*node.nodeServices[service] = append(*node.nodeServices[service], uuid)
		node.serviceList = append(node.serviceList, service)
	}
}

func (node *Node) removeServiceNode(uuid string) {
	for service, nodeList := range node.nodeServices {
		i := slices.Index(*nodeList, uuid)
		fmt.Println(i, *nodeList)
		if i != -1 {
			*nodeList = slices.Delete(*nodeList, i, i+1)
		}
		if len(*nodeList) == 0 {
			delete(node.nodeServices, service)
		}
	}
}
