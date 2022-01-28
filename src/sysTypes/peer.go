package sysTypes

import (
	"sync"
)

type Peer struct {
	*sync.Mutex
	peerList []string
	NumPeers int
}

func NewPeer(peerList []string, numPeers int) Peer {
	return Peer{&sync.Mutex{}, peerList, numPeers}
}

// Adds a peer to the peer list ensuring thread safety
func (p *Peer) Append(peer string) {
	p.Lock()
	defer p.Unlock()
	p.peerList = append(p.peerList, peer)
}

// Gets an individual peer ensuring thread safety
func (p *Peer) GetPeerIndex(index int) string {
	p.Lock()
	defer p.Unlock()
	return p.peerList[index]
}

func (p *Peer) GetPeer() Peer {
	return p.clone()
}

func (p *Peer) clone() Peer {
	return Peer{&sync.Mutex{}, p.peerList, p.NumPeers}
}

func (p *Peer) GetPeerList() []string {
	p.Lock()
	defer p.Unlock()
	return p.peerList
}
