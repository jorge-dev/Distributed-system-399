// Main application entry point
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/jorge-dev/Distributed-system-559/src/common"
	peercommunicator "github.com/jorge-dev/Distributed-system-559/src/peerCommunicator"
	"github.com/jorge-dev/Distributed-system-559/src/registry"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(common.GetlogLevel())
}

var wg sync.WaitGroup

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

	// initialize variables
	// peer := sysTypes.NewPeer(nil, 0)
	// sources := []sysTypes.Source{sysTypes.NewSource(tcpAddr, &peer)}
	peerUdpCommunicator := peercommunicator.NewPeerCommunicator(udpAddr)
	registryClient := registry.NewClient(tcpAddr, udpAddr, teamName)
	peerInfo := sysTypes.NewPeerInfo()
	snips := sysTypes.NewSnips()
	ctx, cancel := context.WithCancel(context.Background())

	// Start the tcp registry
	wg.Add(1)
	go startTCPRegistry(registryClient, peerInfo, snips, ctx)
	waitToExit(ctx, cancel)

	// Start the udp registry
	wg.Add(1)
	go startUdpPeer(peerUdpCommunicator, peerInfo, snips, ctx)
	waitToExit(ctx, cancel)

	log.Info("Program is shutting down")

}

func startTCPRegistry(reg *registry.Client, peerInfo *sysTypes.PeerInfo, snips *sysTypes.Snips, ctx context.Context) {
	defer wg.Done()
	// Start the TCP server
	if err := reg.Start(peerInfo, snips, ctx); err != nil {
		log.Errorf("Error while trying to start the TCP server due to following error: \n %v", err)
	}
}

func startUdpPeer(udpPeer peercommunicator.PeerCommunicator, peerInfo *sysTypes.PeerInfo, snips *sysTypes.Snips, ctx context.Context) {
	defer wg.Done()
	// Start the UDP server
	if err := udpPeer.Start(peerInfo, snips, ctx); err != nil {
		log.Errorf("Error while trying to start the UDP server due to following error: \n %v", err)
	}
}

func waitToExit(ctx context.Context, cancel context.CancelFunc) {
	killReceived := false
	killChannel := make(chan os.Signal, 1)
	signal.Notify(killChannel, os.Interrupt)

	// channel to wait for completion
	waitChannel := make(chan struct{})

	go func() {
		wg.Wait()
		close(waitChannel)
	}()

	select {
	case <-killChannel:
		log.Fatal("Received kill signal, program is shutting down")
		cancel()
		killReceived = true
	case <-waitChannel:
	case <-ctx.Done():
	}

	wg.Wait()

	if killReceived {
		log.Warn("Program being killed by OS signal")
		os.Exit(23)

	}

}
