package protocols

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type Connection interface {
	ConnectToClient()
}

type TCP struct {
	address string
}

func NewTCP(address string) TCP {
	return TCP{address}
}

func (c *TCP) ConnectToClient() (*net.TCPConn, error) {
	// attempt to resolve the address
	addr, err := net.ResolveTCPAddr("tcp", c.address)
	if err != nil {
		log.Debugf("Error while trying to resolve to %s \n", c.address)
		return nil, err
	}

	// attempt to connect to the address
	connection, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Debugf("Error while trying to connect to %s \n", c.address)
		return nil, err
	}

	return connection, nil

}
