package peerProc

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Holds the Peer information
type PeerInfo struct {
	peerAddress   string
	sourceAddress string
	isAlive       bool
	lastSeen      time.Time
}

// Holds the Sent peer Info
type SentPeerInfo struct {
	peerAddr     string
	receiverAddr string
	timestamp    time.Time
}

//Holds the Received Peer Info
type ReceivedPeerinfo struct {
	peerAddrReceived string
	peerAddrSender   string
	timestamp        time.Time
}

// Adds a peer to the list of peers
func AddPeer(peerAddress string, sourceAddress string) {
	// check if the peerAddress is already in the list
	if PeerListIndexLookUp(peerAddress) == -1 {
		fmt.Println("Adding peer with peerAddress " + peerAddress)
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{peerAddress, sourceAddress, true, time.Now()})
		mutex.Unlock()
	} else if PeerListIndexLookUp(sourceAddress) == -1 {
		fmt.Println("Adding peer with SourceAddress" + sourceAddress)
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{sourceAddress, sourceAddress, true, time.Now()})
		mutex.Unlock()

	}

}

// This function formats the list into a string and returns a list of received peers as a string
func PreparelistReceivedPeerInfoToString() string {
	var numPeersRec string = strconv.Itoa(len(listReceivedPeerinfo))
	var peerListRec string = numPeersRec + "\n"
	for _, peer := range listReceivedPeerinfo {
		peerListRec += peer.peerAddrSender + " " + peer.peerAddrReceived + " " + peer.timestamp.Format("2006-01-02 15:04:05") + "\n"
	}
	return peerListRec
}

// This function formats the list into a string and returns a list of sent peers as a string
func PreparelistSentPeerInfoToString() string {
	var numPeersSent string = strconv.Itoa(len(listSentPeerInfo))
	var peerListSent string = numPeersSent + "\n"
	for _, peer := range listSentPeerInfo {
		peerListSent += peer.receiverAddr + " " + peer.peerAddr + " " + peer.timestamp.Format("2006-01-02 15:04:05") + "\n"
	}
	return peerListSent
}

// This function formats the list into a string and returns a list of peers as a string
func PreparelistPeersToString() string {

	var numPeers string = strconv.Itoa(len(listPeers))
	var peerList string = numPeers + "\n"
	for _, peer := range listPeers {
		peerList += peer.peerAddress + "\n"
	}
	return peerList

}

// Handles all peers that have not communicated with the server in the last 5 seconds
func HandleInactivePeers(sourceAddress string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 90):
		}

		mutex.Lock()
		if len(listPeers) > 0 {
			for i := 0; i < len(listPeers); i++ {
				if listPeers[i].peerAddress != sourceAddress {
					if time.Since(listPeers[i].lastSeen) > time.Second*90 {
						listPeers[i].isAlive = false
						fmt.Printf("No communication from peer %v in the last minute and a half.Peer is now inactive\n", listPeers[i].peerAddress)
					}
				}
			}
		}
		mutex.Unlock()
	}
}

// This functions sends a random peer in my list to all the peers in the list
func MulticastMessage(sourceAddress string, conn *net.UDPConn, context context.Context) {

	for {
		rand.Seed(time.Now().UnixNano())
		select {
		case <-context.Done():
			return
		case <-time.After(time.Second * 5):

		}

		if len(listPeers) > 0 {
			peerCount := 0

			// send a random peer to all peers
			peerlen := len(listPeers)
			randPeer := listPeers[rand.Intn(peerlen)]
			for _, peer := range listPeers {
				if !peer.isAlive {
					randPeer = listPeers[rand.Intn(peerlen)]
				}
			}

			for _, peer := range listPeers {
				if CheckForValidAddress(peer.peerAddress) && peer.isAlive {
					sendMessage(peer.peerAddress, UDP_PEER+randPeer.peerAddress, conn)
					mutex.Lock()
					listSentPeerInfo = append(listSentPeerInfo, SentPeerInfo{peer.peerAddress, peer.peerAddress, time.Now()})
					mutex.Unlock()
					peerCount++
				}
			}

		}

	}

}

// Handy function to check it a peerAddress is in the list
func PeerListIndexLookUp(peerAddr string) int {
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].peerAddress == peerAddr {
			return i
		}
	}
	return -1
}

// When a peer is received is stores it in the list of peers
func StorePeers(peerAddr string, senderAddr string) {
	AddPeer(peerAddr, senderAddr)
	mutex.Lock()
	listReceivedPeerinfo = append(listReceivedPeerinfo, ReceivedPeerinfo{peerAddr, senderAddr, time.Now()})
	mutex.Unlock()

}
