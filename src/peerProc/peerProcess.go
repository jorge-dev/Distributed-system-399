package peerProc

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// Define the Registry's request types

	UDP_STOP string = "stop"
	UDP_SNIP string = "snip"
	UDP_PEER string = "peer"
)

var listPeers []PeerInfo
var listSnips []Snip
var listReceivedPeerinfo []ReceivedPeerinfo
var listSentPeerInfo []SentPeerInfo

var mutex = &sync.Mutex{}
var currentTime int = 0
var mainUdpAddress string

// Handles the UDP responsability concurrently
func PeerProcess(conn *net.UDPConn, sourceAddress string, ctx context.Context) {
	mainUdpAddress = sourceAddress
	// listPeers = append(listPeers, PeerInfo{sourceAddress, sourceAddress, true, time.Now()})
	fmt.Printf("Peer Party Started at %s\n", sourceAddress)
	wg := sync.WaitGroup{}
	childCtx, cancel := context.WithCancel(ctx)
	wg.Add(4)
	go func() {
		defer wg.Done()
		messageHandler(conn, sourceAddress, childCtx, cancel)
	}()

	go func() {
		defer wg.Done()
		SnipHandler(sourceAddress, conn, childCtx)
	}()

	go func() {
		defer wg.Done()
		MulticastMessage(sourceAddress, conn, childCtx)
	}()

	go func() {
		defer wg.Done()
		HandleInactivePeers(sourceAddress, childCtx)
	}()
	wg.Wait()

}

// Connects to the UDp server and sends a message to the specified address
func sendMessage(peerAddress, msg string, conn *net.UDPConn) {
	udpAdd, err := net.ResolveUDPAddr("udp", peerAddress)
	if err != nil {
		fmt.Println("Error in resolving UDP address, error is: ", err)
		return
	}

	_, err = conn.WriteToUDP([]byte(msg), udpAdd)
	if err != nil {
		fmt.Printf("Error while sending message to %s due to following error: \n %v", peerAddress, err)
		return
	}
	// fmt.Printf("Message sent to %s: %s\n", peerAddress, msg)

}

// TThis function will try to check if an address is valid by trying to get a response
func CheckForValidAddress(address string) bool {
	// check if the host and port are valid

	if _, err := net.ResolveUDPAddr("udp", address); err != nil {
		return false
	}
	return true
}

// This function handles the UDP messages commands
func messageHandler(conn *net.UDPConn, sourceAddress string, ctx context.Context, cancel context.CancelFunc) {

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Closing the connection")
			return
		default:
			msg, senderAddr, err := receiveUdpMessage(sourceAddress, conn)
			if err != nil {
				fmt.Println("Error while receiving message: ", err)
				continue
			}
			// update last seen
			for i := 0; i < len(listPeers); i++ {
				if listPeers[i].peerAddress == senderAddr {
					listPeers[i].lastSeen = time.Now()
					listPeers[i].isAlive = true
				}
				// AddPeer(senderAddr, sourceAddress)

			}
			// only focus on first 4 characters
			// fmt.Pcrintf("Message received from %s: %s\n", senderAddr, msg)
			if len(msg) >= 4 {
				switch msg[:4] {
				case UDP_STOP:
					fmt.Println("Stopping UDP server")
					conn.Close()
					cancel()
					return
				case UDP_SNIP:
					fmt.Println("Receiving Snips")
					command := strings.Trim(msg[4:], "\n")
					go storeSnips(command, senderAddr)
				case UDP_PEER:
					// fmt.Println("Peer info received")
					peerAddr := strings.Trim(msg[4:], "\n")
					go StorePeers(peerAddr, senderAddr)
				default:
					fmt.Printf("Unknown command received from %s: %s\n", senderAddr, msg)

				}
			} else {
				fmt.Println("Message is not long enough to be a command")
			}
		}
	}
}

// Gets the max value from two values
func getMAxValue(val1, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}

// Handles messages received from other peers
func receiveUdpMessage(address string, conn *net.UDPConn) (string, string, error) {

	// Read from the connection
	data := make([]byte, 1024)
	len, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return "", "", err
	}
	msg := strings.TrimSpace(string(data[:len]))

	return msg, addr.String(), nil

}
