package main

import (
	"fmt"
	"go_server/client"
	"go_server/server"
	"go_server/utils"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

var (
	mode        = utils.Node
	local       string
	localValid  = true
	remote      string
	remoteValid = true
)

func main() {
	ParseArgs()

	switch mode {
	case utils.Node:
		server.HandleNode(local, remote)
	case utils.Client:
		client.HandleClient(local, remote)
	}
}

func ParseArgs() {
	// Parse flags
	var modeArg string
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "==Usage: main.go [options]\n")
		fmt.Fprintln(os.Stderr, "Options:")
		pflag.PrintDefaults()
	}
	pflag.StringVarP(&modeArg, "mode", "m", "", "server|client (required)")
	pflag.StringVarP(&local, "local", "l", "", "Local IP address and port (e.g. localhost:8080 | :8080)")
	pflag.StringVarP(&remote, "remote", "r", "", "Remote IP address and port (e.g. 8.8.8.8:8080)")
	pflag.Parse()

	// Validate mode
	valid := true
	modeArg = strings.ToLower(modeArg)
	if _, v := utils.EntityTypes[modeArg]; !v {
		fmt.Printf("==ERROR: Invalid MODE\n  Provided: %s, expected 'Balancer', 'Server', or 'Client'\n", modeArg)
		valid = false
	}
	mode = utils.EntityTypes[modeArg]
	if mode == utils.Leader {
		valid = false
	}

	var errorMessage []string
	var remotePort int
	var localPort int
	var err error
	// Validate host
	switch mode {
	case utils.Node:
		// Process remote (optional)
		if remote != "" {
			// Validate remote
			if !strings.Contains(remote, ":") {
				errorMessage = append(errorMessage, "Port not provided, using default: 8080")
				errorMessage = append(errorMessage, "  INFO: Port input must use a colon (e.g. example_address:8080)")
				remotePort = 8080
				remote = fmt.Sprintf("%s:%d", remote, remotePort)

				// Attempt a connection to ensure the remote is proper
				listen, err := net.Dial(string(utils.TCP), remote)
				if err != nil {
					remoteValid = false
					errorMessage = append(errorMessage, "Could not contact remote:")
					errorMessage = append(errorMessage, fmt.Sprintf("  Ensure remote is correct. Expected: example:port, received %s", remote))
					goto NODE_LOCAL
				}
				listen.Close()
				goto NODE_LOCAL
			}

			remoteValidation := strings.Split(remote, ":")
			if len(remoteValidation) != 2 {
				remoteValid = false
				errorMessage = append(errorMessage, "Incorrect remote format provided\n")
				errorMessage = append(errorMessage, fmt.Sprintf("  Expected: example_address:port, received %s\n", remote))
				goto NODE_LOCAL
			}

			if remoteValidation[0] == "" {
				remoteValid = false
				errorMessage = append(errorMessage, "Invalid input for remote\n")
				errorMessage = append(errorMessage, fmt.Sprintf("  Expected: example_address:port, received %s\n", remote))
				goto NODE_LOCAL
			}

			remotePort, err = strconv.Atoi(remoteValidation[1])
			if err != nil {
				errorMessage = append(errorMessage, "Could not parse the port provided, using default: 8080\n")
				remotePort = 8080
			}
			remote = fmt.Sprintf("%s:%d", remoteValidation[0], remotePort)
		}

		// Process local (optional)
	NODE_LOCAL:
		if local != "" {
			if !strings.Contains(local, ":") {
				errorMessage = append(errorMessage, "INFO: Port input must use a colon (e.g. example_address:8080)")
				errorMessage = append(errorMessage, "  Port not provided, using default: 8080")
				localPort = 8080
				local = fmt.Sprintf("%s:%d", local, localPort)

				// Attempt a connection to ensure the remote is proper
				listen, err := net.Listen(string(utils.TCP), local)
				if err != nil {
					localValid = false
					errorMessage = append(errorMessage, "Incorrect remote format provided")
					errorMessage = append(errorMessage, fmt.Sprintf("  Expected: example_address:port, received %s", local))
					goto NODE_END
				}
				listen.Close()
				goto NODE_END
			}

			if len(strings.Split(local, ":")) != 2 {
				localValid = false
				errorMessage = append(errorMessage, "Incorrect Local format provided")
				errorMessage = append(errorMessage, fmt.Sprintf("  Expected: example_address:port, received %s", local))
				goto NODE_END
			}

			localValidation := strings.Split(local, ":")
			if localValidation[0] == "" {
				fmt.Println("No Local information provided, using default: 0.0.0.0")
				localValidation[0] = "0.0.0.0"
			}

			localPort, err = strconv.Atoi(localValidation[1])
			if err != nil {
				fmt.Println("Could not parse the port provided, attempting with default port: 8080")
				localPort = 8080
			}
			local = fmt.Sprintf("%s:%d", localValidation[0], localPort)
			goto NODE_END
		}

		fmt.Println("No Local information provided, using defaults.")
		local = fmt.Sprintf("%s:%d", "0.0.0.0", 8080)

	NODE_END:
		if !localValid || !remoteValid {
			valid = false
		}

	case utils.Client:
		if remote == "" {
			errorMessage = append(errorMessage, "Invalid input:\n")
			errorMessage = append(errorMessage, "  Expected: example_address:port, nothing was recieved\n")
			remoteValid = false
			valid = false
		} else {

		}
		fmt.Println("Not set up yet kekekekekeke")
		// Process remote (required)
	}

	for _, msg := range errorMessage {
		fmt.Println(msg)
	}
	if !valid {
		fmt.Printf("Could not initialize %s.\n", modeArg)
		pflag.Usage()
		os.Exit(1)
	}

	fmt.Println(modeArg, "starting on", local)
}
