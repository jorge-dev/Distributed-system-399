package protocols

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type UDP struct {
	address string
}

func NewDUP(address string) UDP {
	return UDP{address}
}

func (u *UDP) ConnectToClient() (*net.UDPConn, error) {
	// attempt to resolve the address
	log.Infof("Initializing UDP server on: %s", u.address)
	udpAddr, err := net.ResolveUDPAddr("udp", u.address)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s\n", u.address)
		return nil, err
	}
	// checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	// checkError(err)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s\n", u.address)
		return nil, err
	}

	return conn, nil

}
