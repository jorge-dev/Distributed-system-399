package protocols

import "net"

type Connection interface {
	ConnectToClient()
	ListenToClient()
}

type TCP struct {
	address string
}

func NewTCP(address string) TCP {
	return TCP{address}
}

func (c *TCP) ConnectToClient() (net.Conn, error) {
	// attempt to resolve the address
	addr, err := net.ResolveTCPAddr("tcp", c.address)
	if err != nil {
		return nil, err
	}

	// attempt to connect to the address
	connection, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	return connection, nil

}

func (c *TCP) ListenToClient() {

}
