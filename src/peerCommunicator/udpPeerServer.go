package peercommunicator

import (
	"context"
	"net"
	"sync"

	"github.com/jorge-dev/Distributed-system-559/src/protocols"
	log "github.com/sirupsen/logrus"
)

type PeerCommunicator struct {
	udpAddr    string
	connection *net.UDPConn
}

func NewPeerCommunicator(udpAddr string) PeerCommunicator {
	return PeerCommunicator{udpAddr: udpAddr}
}

var wg sync.WaitGroup

func (pc *PeerCommunicator) Start(ctx context.Context) error {

	udpProto := protocols.NewDUP(pc.udpAddr)
	var err error
	pc.connection, err = udpProto.ConnectToClient()
	if err != nil {
		log.Debugf("Error while trying to connect to %s due to following error: \n %v", pc.udpAddr, err)
		return err

	}
	defer pc.connection.Close()

	return nil
}
