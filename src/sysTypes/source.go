package sysTypes

import (
	"github.com/jorge-dev/Distributed-system-559/src/common"
	// "github.com/jorge-dev/Distributed-system-559/sysTypes/peer"
)

// import (
// 	"github.com/jorge-dev/Distributed-system-559/sysTypes"
// )

type Source struct {
	address   string
	peer      *Peer
	timeStamp string
}

func (source *Source) GetSourceAddress() string {
	return source.address
}

func NewSource(address string, peers *Peer) Source {
	return Source{address, peers, common.GetCurrentDateTime()}
}

func (source *Source) GetTimeStamp() string {
	return source.timeStamp
}

func (source *Source) UpdateTimeStamp() {
	source.timeStamp = common.GetCurrentDateTime()
}

func (source *Source) SetAddress(address string) {
	source.address = address
}

func (source *Source) GetAddress() string {
	return source.address
}

func (source *Source) GetPeerType() Peer {
	return source.peer.GetPeer()
}
