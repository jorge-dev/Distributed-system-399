// Main application entry point
package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/jorge-dev/Distributed-system-559/src/client"
)

func main() {
	rand.Seed(time.Now().UnixNano())
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

	var teamName string = "Jorge Avila" + strconv.Itoa(rand.Intn(1000))
	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(2)

	// connect to the server
	go func() {
		defer wg.Done()
		err := client.ConnectTCP(tcpHost, tcpPort, udpHost, udpPort, teamName, ctx)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}()

	// connect to the Udp server
	go func() {
		defer wg.Done()
		err := client.ConnectUdpServer(teamName, udpHost, udpPort, ctx)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("UDP server started")
		err2 := client.ConnectTCP(tcpHost, tcpPort, udpHost, udpPort, teamName, ctx)
		if err2 != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}

		cancel()
	}()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(osSignal, os.Interrupt, os.Kill)

	select {
	case <-osSignal:
		fmt.Println("\nShutting down gracefully ...")
		cancel()

	case <-ctx.Done():
		fmt.Println("\nShutting down gracefully...")
	}
	wg.Wait()

}
