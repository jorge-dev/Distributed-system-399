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

func ReceiveMessage(conn *net.UDPConn) (string, string, error) {
	buffer := make([]byte, 1024)
	len, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Debugf("Error while receiving message due to following error: \n %v", err)
		return "", "", err
	}
	log.WithField("data", string(buffer[0:len])).Debugf("Message received from %s", addr)
	receivedMessage := string(buffer[0:len])
	return receivedMessage, addr.String(), nil
}

func SendMessage(conn *net.UDPConn, peerAddress, data string) {
	udpAdd, err := net.ResolveUDPAddr("udp", peerAddress)
	if err != nil {
		log.Errorf("Error in resolving UDP address, error is: ", err)
		return
	}

	_, err = conn.WriteToUDP([]byte(data), udpAdd)
	if err != nil {
		log.Errorf("Error while sending message to %s due to following error: \n %v", peerAddress, err)
		return
	}
	log.WithField("data", data).Debugf("Message sent to %s", peerAddress)
	// fmt.Printf("Message sent to %s: %s\n", peerAddress, msg)
}

func CheckForValidAddress(address string) bool {
	_, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s\n", address)
		return false
	}
	return true
}
