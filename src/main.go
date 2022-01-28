package main

import (
	"fmt"
	"os"

	"github.com/jorge-dev/Distributed-system-559/src/client"
)

func main() {
	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: client host port")
		os.Exit(1)

	}

	host := args[0]
	port := args[1]

	// connect to the server
	err := client.Connect(host, port)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
