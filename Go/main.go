package main

import (
	"fmt"
	"go_server/balancer"
	"go_server/client"
	"go_server/node"
	"go_server/utils"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

var (
	host string
	mode = utils.Server
)

func main() {
	ParseArgs()

	switch mode {
	case utils.Server:
		node.HandleNode(host)
	case utils.Client:
		client.HandleClient(host)
	case utils.Balancer:
		balancer.Handlebalancer(host)
	}
}

func ParseArgs() {
	// Parse flags
	var modeArg string
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: main.go [options]\n")
		fmt.Fprintln(os.Stderr, "Options:")
		pflag.PrintDefaults()
	}
	pflag.StringVarP(&host, "host", "h", "", "hostname:port (required)")
	pflag.StringVarP(&modeArg, "mode", "m", "", "balancer|server|client (required)")
	pflag.Parse()

	// Validate host
	valid := true
	if host == "" || !strings.Contains(host, ":") || strings.Count(host, ":") > 1 {
		fmt.Printf("==ERROR: Invalid HOST\n  Passed: %s, expected IP address or Domain.\n", host)
		valid = false
	}

	// Ensure prt is passed correctly
	prt, err := strconv.Atoi(strings.Split(host, ":")[1])
	if err != nil || prt < 2000 || prt > 32000 {
		fmt.Printf("==ERROR: Invalid PORT\n  Provided: %s, expected value between 2000 and 32000\n", strings.Split(host, ":")[1])
		valid = false
	}

	// Validate mode
	modeArg = strings.ToLower(modeArg)
	if _, valid := utils.ValidModes[modeArg]; !valid {
		fmt.Printf("==ERROR: Invalid MODE\n  Provided: %s, expected 'Balancer', 'Server', or 'Client'\n", modeArg)
		valid = false
	}
	mode = utils.ValidModes[modeArg]

	if !valid {
		pflag.Usage()
		os.Exit(1)
	}

	fmt.Println(modeArg, "starting on", host)
}
