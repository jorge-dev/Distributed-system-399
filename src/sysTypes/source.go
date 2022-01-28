// Declare thr Source type (Data Structure) and its methods

package sysTypes

import (
	"github.com/jorge-dev/Distributed-system-559/src/common"
)

// Struct used to store the source information
type Source struct {
	address   string
	peer      *Peer
	timeStamp string
}

// Parameterized constructor for Source (Go's version of a constructor)
func NewSource(address string, peers *Peer) Source {
	return Source{address, peers, common.GetCurrentDateTime()}
}

func (source *Source) GetSourceAddress() string {
	return source.address
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
