package sysTypes

import (
	"context"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/jorge-dev/Distributed-system-559/src/protocols"
	log "github.com/sirupsen/logrus"
)

type MainInfo struct {
	PeerAddr   string
	SourceAddr string
	IsAlive    bool
	LastSeen   time.Time
}

type SentInfo struct {
	SendToAddr   string
	PeerSentAddr string
	Timestamp    time.Time
}

type ReceivedInfo struct {
	ReceivedPeerAddr string
	SenderAddr       string
	Timestamp        time.Time
}

type PeerInfo struct {
	MainPeerInfo     []MainInfo
	SentPeerInfo     []SentInfo
	receivedPeerInfo []ReceivedInfo
}

func NewPeerInfo() *PeerInfo {
	return &PeerInfo{}
}

// Adds a peer to the main info peer list ensuring thread safety access
func (p *PeerInfo) AppendMainInfoPeer(peerAddr string, sourceAddr string) {
	mutex := &sync.Mutex{}
	for _, peer := range p.MainPeerInfo {
		if peer.PeerAddr == peerAddr {
			return
		} else if peer.SourceAddr == sourceAddr {
			return
		} else if peer.PeerAddr == sourceAddr {
			return
		}
	}
	mutex.Lock()

	p.MainPeerInfo = append(p.MainPeerInfo, MainInfo{PeerAddr: peerAddr, SourceAddr: sourceAddr, IsAlive: true, LastSeen: time.Now()})
	mutex.Unlock()
}

// Adds a peer to the sent info peer list ensuring thread safety access
func (p *PeerInfo) AppendSentInfoPeer(sendToAddress string, peerSentAddress string) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	p.SentPeerInfo = append(p.SentPeerInfo, SentInfo{SendToAddr: sendToAddress, PeerSentAddr: peerSentAddress, Timestamp: time.Now()})
	mutex.Unlock()
}

// Adds a peer to the received info peer list ensuring thread safety access
func (p *PeerInfo) AppendReceivedInfoPeer(receivedPeerAddress string, senderAddress string) {
	// get the index of receivedPeerAddress and senderAddress in MainInfo
	log.Debugf("In appedndReceivedInfoPeer")
	mutex := &sync.Mutex{}
	receivedPeerIndex := p.getMainInfoPeerIndex(receivedPeerAddress)
	senderIndex := p.getMainInfoSourceIndex(senderAddress)
	isValidAddress := protocols.CheckForValidAddress(receivedPeerAddress)

	// if the peer is not in the list of MainInfo peers, add it
	if receivedPeerIndex == -1 && senderIndex == -1 && isValidAddress {
		p.AppendMainInfoPeer(receivedPeerAddress, senderAddress)
	} else if receivedPeerIndex == -1 && senderIndex != -1 && isValidAddress {
		p.AppendMainInfoPeer(receivedPeerAddress, receivedPeerAddress)
	} else if receivedPeerIndex != -1 && senderIndex == -1 && isValidAddress {
		p.AppendMainInfoPeer(senderAddress, senderAddress)
	}
	log.Debugf("This is the current list og mainInfo peers: %v", p.MainPeerInfo)
	mutex.Lock()
	p.receivedPeerInfo = append(p.receivedPeerInfo, ReceivedInfo{ReceivedPeerAddr: receivedPeerAddress, SenderAddr: senderAddress, Timestamp: time.Now()})
	// log.Debugf("Added to received peer info list: %v", p.receivedPeerInfo)

	mutex.Unlock()

}

// This function formats the main info list into a string and returns a list as a string
func (p *PeerInfo) GetPeerMainInfoAsString() string {
	totalPeers := strconv.Itoa(len(p.MainPeerInfo))
	peerList := totalPeers + "\n"
	for _, peer := range p.MainPeerInfo {
		peerList += peer.PeerAddr + "\n"
	}
	return peerList
}

// update main info last seen time and alive status
func (p *PeerInfo) UpdateMainInfoLastSeen(peerAddr string) {
	for i, peer := range p.MainPeerInfo {
		if peer.PeerAddr == peerAddr {
			p.MainPeerInfo[i].LastSeen = time.Now()
			p.MainPeerInfo[i].IsAlive = true
		}

	}
}

// This function formats the list into a string and returns a list of sent peers as a string
func (p *PeerInfo) GetSentPeerInfoAsString() string {
	totalPeers := strconv.Itoa(len(p.SentPeerInfo))
	peerList := totalPeers + "\n"
	for _, peer := range p.SentPeerInfo {
		peerList += peer.SendToAddr + " " + peer.PeerSentAddr + "\n" + peer.Timestamp.Format("2006-01-02 15:04:05") + "\n"
	}
	return peerList
}

// This function formats the list into a string and returns a list of received peers as a string
func (p *PeerInfo) GetreceivedPeerInfoAsString() string {
	totalPeers := strconv.Itoa(len(p.receivedPeerInfo))
	peerList := totalPeers + "\n"
	for _, peer := range p.receivedPeerInfo {
		peerList += peer.SenderAddr + " " + peer.ReceivedPeerAddr + "\n" + peer.Timestamp.Format("2006-01-02 15:04:05") + "\n"
	}
	return peerList
}

// Handles peers that have not communicated with the server for the last 10 seconds
func (p *PeerInfo) HandleDeadPeers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("PeerInfo: handleDeadPeers: Done")
			return
		case <-time.After(10 * time.Second):
		}
		// p.Lock()
		// defer p.Unlock()
		if len(p.MainPeerInfo) > 0 {
			for i, peer := range p.MainPeerInfo {
				if time.Since(peer.LastSeen) > 10*time.Second {
					p.MainPeerInfo[i].IsAlive = false
					log.Debugf("PeerInfo: handleDeadPeers: Peer %s is dead", peer.PeerAddr)
				}
			}
		}
	}
}

// Sends a heartbeat to all peers in the MainInfo peer list every 5 seconds
func (p *PeerInfo) SendHeartbeat(connection *net.UDPConn, ctx context.Context) {
	// seed a random number generator

	for {
		select {
		case <-ctx.Done():
			log.Debug("PeerInfo: sendHeartbeat: Done")
			return
		case <-time.After(5 * time.Second):
		}

		if len(p.MainPeerInfo) > 0 {
			peerSentCount := 0
			// get a random peer from the MainInfo peer list
			randPeer := p.getRandomPeer()
			for _, peer := range p.MainPeerInfo {
				if peer.IsAlive && protocols.CheckForValidAddress(peer.PeerAddr) {
					// send a heartbeat to all peers
					log.Debugf("PeerInfo: sendHeartbeat: Sending heartbeat to %s", peer.PeerAddr)
					data := "peer" + peer.PeerAddr
					go protocols.SendMessage(connection, peer.PeerAddr, data)
					p.AppendSentInfoPeer(peer.PeerAddr, randPeer.PeerAddr)
					peerSentCount++
				}
			}
		}
	}

}

// get a random peer from the main info peer list if its alive
func (p *PeerInfo) getRandomPeer() MainInfo {
	mutex := &sync.Mutex{}
	rand.Seed(time.Now().UnixNano())
	// select a random peer from the MainInfo peer list
	randPeerIndex := rand.Intn(len(p.MainPeerInfo))
	randPeer := p.MainPeerInfo[randPeerIndex]
	// check if the random peer is alive, if not regenerate a new random peer
	mutex.Lock()

	for !randPeer.IsAlive {
		randPeerIndex = rand.Intn(len(p.MainPeerInfo))
		randPeer = p.MainPeerInfo[randPeerIndex]
	}
	mutex.Unlock()
	return randPeer
}

// Returns the main info peer list
func (p *PeerInfo) GetMainInfoPeerList() []MainInfo {
	return p.MainPeerInfo
}

// Returns the sent info peer list
func (p *PeerInfo) GetSentInfoPeerList() []SentInfo {
	return p.SentPeerInfo
}

// Returns the received info peer list
func (p *PeerInfo) GetReceivedInfoPeerList() []ReceivedInfo {
	return p.receivedPeerInfo
}

// get index of peer address in MainInfo
func (p *PeerInfo) getMainInfoPeerIndex(peerAddr string) int {
	for i, peer := range p.MainPeerInfo {
		if peer.PeerAddr == peerAddr {
			return i
		}
	}
	return -1
}

// get index of peer source address in MainInfo
func (p *PeerInfo) getMainInfoSourceIndex(sourceAddr string) int {
	for i, peer := range p.MainPeerInfo {
		if peer.SourceAddr == sourceAddr {
			return i
		}
	}
	return -1
}

// Check if peer address is in MainInfo
func (p *PeerInfo) checkMainInfoPeer(peerAddr string) bool {
	for _, peer := range p.MainPeerInfo {
		if peer.PeerAddr == peerAddr {
			return true
		}
	}
	return false
}

// check if source address is in MainInfo
func (p *PeerInfo) checkMainInfoSource(sourceAddr string) bool {
	for _, peer := range p.MainPeerInfo {
		if peer.SourceAddr == sourceAddr {
			return true
		}
	}
	return false
}
