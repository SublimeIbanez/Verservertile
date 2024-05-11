package server

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func (node *Node) outputHandler() {
	for out := range node.outputChannel {
		select {
		case <-node.context.Done():
			return

		default:
			fmt.Println(out)
		}
	}
}

func (node *Node) inputHandler(wait *sync.WaitGroup) {
	defer wait.Done()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		select {
		case <-node.context.Done():
			return

		default:
			text := strings.ToLower(scanner.Text())
			node.outputChannel <- "You typed " + text
			if text == "exit" {
				node.context.Done()
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		node.outputChannel <- "Error reading from the console: " + err.Error()
	}
}

func (node *Node) printNodeList() {
	node.outputChannel <- "Node List:"
	for id, a := range node.nodes {
		node.outputChannel <- fmt.Sprintf("  - %s::%s", id, a)
	}
	node.outputChannel <- ""
}

func (node *Node) printNodeServiceList() {
	node.outputChannel <- "Node Services List:"
	for service, nodeList := range node.nodeServices {
		node.outputChannel <- fmt.Sprintf("  Service: %s", service)
		for _, n := range *nodeList {
			node.outputChannel <- fmt.Sprintf("    - %s", n)
		}
	}
	node.outputChannel <- ""
}

func (node *Node) printServiceList() {
	node.outputChannel <- "Services List:"
	for _, service := range node.serviceList {
		node.outputChannel <- fmt.Sprintf("  - %s", service)
	}
	node.outputChannel <- ""
}
