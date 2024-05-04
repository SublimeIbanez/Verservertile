package main

import (
	"fmt"
	"go_server/client"
	"go_server/server"
	"go_server/utils"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

var (
	host string
	mode = utils.Server
	port uint16
)

func main() {
	ParseArgs()

	switch mode {
	case utils.Balancer:
		server.HandleBalancer(port)
	case utils.Server:
		server.HandleNode(host, port)
	case utils.Client:
		client.HandleClient(host)
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
	pflag.StringVarP(&host, "host", "h", "", "host")
	pflag.Uint16VarP(&port, "port", "p", 8080, "port")
	pflag.StringVarP(&modeArg, "mode", "m", "", "balancer|server|client (required)")
	pflag.Parse()

	// Validate mode
	valid := true
	modeArg = strings.ToLower(modeArg)
	if _, v := utils.ValidModes[modeArg]; !v {
		fmt.Printf("==ERROR: Invalid MODE\n  Provided: %s, expected 'Balancer', 'Server', or 'Client'\n", modeArg)
		valid = false
	}
	mode = utils.ValidModes[modeArg]

	// Validate host
	if mode != utils.Balancer && host == "" {
		fmt.Print("==ERROR: Invalid HOST\n  Nothing was provided.\n")
		valid = false
	}
	// Ensure prt is passed correctly
	if port != 8080 && (port < 2000 || port > 32000) {
		fmt.Printf("==ERROR: Invalid PORT\n  Provided: %d, expected value between 2000 and 32000\n", port)
		valid = false
	}

	if !valid {
		pflag.Usage()
		os.Exit(1)
	}

	fmt.Println(modeArg, "starting on", host)
}
