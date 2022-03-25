// Main application entry point
// Jorge Avila
// 10123968
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jorge-dev/Distributed-system-559/src/client"
	"github.com/jorge-dev/Distributed-system-559/src/common"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(common.GetlogLevel())
}

func main() {
	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <udpPort> <flag>")
		os.Exit(1)

	}
	// get commandline args
	udpPort := args[0]
	flag := args[1]

	// initialize the tcp, udp address and team name
	tcpAddr, udpAddr, teamName := common.GetIpAndTeam(flag, udpPort)
	log.WithField("tcpAddr", tcpAddr).WithField("udpAddr", udpAddr).WithField("teamName", teamName).Info("Starting client")

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(2)

	// connect to the server
	go func() {
		defer wg.Done()
		err := client.ConnectTCP(tcpAddr, udpAddr, teamName, ctx)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}()

	// connect to the Udp server
	go func() {
		defer wg.Done()
		err := client.ConnectUdpServer(teamName, udpAddr, ctx)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
		fmt.Println("UDP server started")
		err2 := client.ConnectTCP(tcpAddr, udpAddr, teamName, ctx)
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
