// Main application entry point
package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jorge-dev/Distributed-system-559/src/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Error loading .env file")
		os.Exit(1)
	}

	// fmt.Printf("%s\n", os.Getenv("TEST_ENV"))
	// os.Exit(0)

	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <host> <port>")
		os.Exit(1)

	}

	host := args[0]
	port := args[1]

	// connect to the server
	err = client.Connect(host, port)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}
