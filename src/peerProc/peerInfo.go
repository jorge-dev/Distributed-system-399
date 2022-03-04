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
	// check if the peer is already in the list
	mutex.Lock()
	for _, peer := range listPeers {
		if peer.peerAddress == peerAddress {
			mutex.Unlock()
			return
		} else if peer.sourceAddress == peerAddress {
			mutex.Unlock()
			return
		} else if peer.sourceAddress == sourceAddress {
			mutex.Unlock()
			return
		}
	}
	listPeers = append(listPeers, PeerInfo{peerAddress, sourceAddress, true, time.Now()})
	mutex.Unlock()

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
	for _, peer := range listPeers {
		fmt.Println(peer.peerAddress + "\n" + peer.sourceAddress + "\n" + peer.lastSeen.Format("2006-01-02 15:04:05") + "\n")
	}

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
		case <-time.After(time.Second * 15):
		}

		mutex.Lock()
		// fmt.Printf("Peers in the list: %v\n", listPeers)
		if len(listPeers) > 0 {
			for i := 0; i < len(listPeers); i++ {
				if listPeers[i].peerAddress != sourceAddress {
					if time.Since(listPeers[i].lastSeen) > time.Second*10 {
						listPeers[i].isAlive = false
					}
				}
			}
			fmt.Printf("Inactive peers removed. Peers left %d\n", len(listPeers))
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

		mutex.Lock()
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

			// fmt.Println("Sending peers")
			for _, peer := range listPeers {
				if CheckForValidAddress(peer.peerAddress) && peer.isAlive {
					sendMessage(peer.peerAddress, UDP_PEER+randPeer.peerAddress, conn)
					listSentPeerInfo = append(listSentPeerInfo, SentPeerInfo{peer.peerAddress, peer.peerAddress, time.Now()})
					peerCount++
				}
			}

		}
		mutex.Unlock()
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

// Handy function to check it a peerSource is in the list
func PeerSourceListIndexLookUp(peerSrc string) int {
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].sourceAddress == peerSrc {
			return i
		}

	}
	return -1

}

// When a peer is received is stores it in the list of peers
func StorePeers(peerAddr string, senderAddr string) {
	// Get the peer and source index
	peerIndex := PeerListIndexLookUp(peerAddr)
	sourceIndex := PeerSourceListIndexLookUp(senderAddr)

	// If the peer is not in the list, add it
	if peerIndex == -1 && sourceIndex == -1 && CheckForValidAddress(peerAddr) {
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{peerAddr, senderAddr, true, time.Now()})
		mutex.Unlock()
		fmt.Printf("New peer added: %s\n", peerAddr)
	} else if peerIndex == -1 && sourceIndex != -1 && CheckForValidAddress(peerAddr) {
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{peerAddr, peerAddr, true, time.Now()})
		mutex.Unlock()
		fmt.Printf("New peer added: %s\n", peerAddr)
	} else if peerIndex != -1 && sourceIndex == -1 && CheckForValidAddress(peerAddr) {
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{senderAddr, senderAddr, true, time.Now()})
	}

	listReceivedPeerinfo = append(listReceivedPeerinfo, ReceivedPeerinfo{peerAddr, senderAddr, time.Now()})

}
