// Main application entry point
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/jorge-dev/Distributed-system-559/src/client"
)

func main() {
	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 4 {
		fmt.Println("Usage: go run main.go <tcpHost> <tcpPort> <udpHost> <udpPort>")
		os.Exit(1)

	}

	tcpHost := args[0]
	tcpPort := args[1]

	udpHost := args[2]
	udpPort := args[3]

	wg := sync.WaitGroup{}
	wg.Add(2)

	// connect to the server
	go func() {
		defer wg.Done()
		err := client.ConnectTCP(tcpHost, tcpPort, udpHost, udpPort)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}()

	// connect to the server
	go func() {
		defer wg.Done()
		err := client.ConnectUdpServer(udpHost, udpPort)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}()

	wg.Wait()

}
