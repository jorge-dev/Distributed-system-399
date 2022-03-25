package peerProc

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// Holds the Snip Information
type Snip struct {
	message    string
	senderAddr string
	timeStamp  int
}

// This function formats the list into a string and returns a list of snips received as a string
func PreparelistSnipsToString() string {
	var numSnips string = strconv.Itoa(len(listSnips))
	var snipList string = numSnips + "\n"

	for _, snip := range listSnips {
		snipList += strconv.Itoa(snip.timeStamp) + " " + snip.message + " " + snip.senderAddr + "\n"
	}
	return snipList
}

// Handles what happends when you get a snip
func SnipHandler(sourceAddress string, conn *net.UDPConn, ctx context.Context) {
	ch := make(chan string, 5)
	// var ch chan string
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fmt.Println("Send snip: ", scanner.Text())
			ch <- scanner.Text()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:

			sendSnip(msg, sourceAddress, conn)
		}
	}
}

//Send a snip to the peer
func sendSnip(msg string, sourceAddress string, conn *net.UDPConn) {
	mutex.Lock()
	currentTime++
	snipCurrentTime := strconv.Itoa(currentTime)
	msg = "snip" + snipCurrentTime + " " + msg
	for _, peer := range listPeers {
		if CheckForValidAddress(peer.peerAddress) {
			if peer.isAlive {
				sendMessage(peer.peerAddress, msg, conn)
			}
		} else {
			fmt.Printf("Invalid address: %s\n", peer.peerAddress)
		}
	}
	mutex.Unlock()
}

// After receiving a snip store it for report
func storeSnips(command string, senderAddr string) {
	msg := strings.Split(command, " ")
	timestamp, err := strconv.Atoi(msg[0])
	if err != nil {
		fmt.Println("Timestamp is not a valid number")
		return
	}
	if len(msg) < 2 {
		fmt.Printf("Invalid snip command: \n message: %s%s\n", command, msg)
		return
	}
	// Store the snip to list
	// join the rest of the message
	snipContent := strings.Join(msg[1:], " ")

	mutex.Lock()
	// check which time is the latest
	if currentTime != timestamp {
		currentTime = getMAxValue(currentTime, timestamp)
	}

	listSnips = append(listSnips, Snip{snipContent, senderAddr, currentTime})
	mutex.Unlock()

	// update last seen
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].peerAddress == senderAddr {
			listPeers[i].lastSeen = time.Now()
			listPeers[i].isAlive = true
		}
	}

}
