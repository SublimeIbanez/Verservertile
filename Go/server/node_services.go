package server

import (
	"encoding/json"
	"fmt"
	"go_server/common"
	"net"
	"slices"
)

type Service string

const (
	Chat Service = "Chat"
)

var ServiceList = []Service{
	Chat,
}

func (node *Node) initServices() {
	node.updateServices(ServiceList, node.uuid)
}

func (node *Node) updateServices(services []Service, uuid string) {
	for _, s := range services {
		if node.nodeServices[s] == nil {
			node.nodeServices[s] = &[]string{}
		}
		*node.nodeServices[s] = append(*node.nodeServices[s], uuid)
		node.serviceList = append(node.serviceList, s)
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

func (node *Node) handleService(base *common.BaseMessage, conn *net.Conn) {
	var request common.ServiceOp
	err := json.Unmarshal(*base.Data[common.ServiceOperation], &request)
	if err != nil {
		node.outputChannel <- "Could not unmarshal service operation request"
		(*conn).Close()
		return
	}

	switch Service(request.Service) {
	case Chat:
		{

		}
	}
}
