package sysTypes

import (
	"bufio"
	"context"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/jorge-dev/Distributed-system-559/src/protocols"
	log "github.com/sirupsen/logrus"
)

type singleSnip struct {
	message    string
	senderAddr string
	timeStamp  int
}
type Snips struct {
	snips []singleSnip
}

func NewSnips() *Snips {
	return &Snips{}
}

var currentTimeStamp int = 0

// Appends a new snip to the snips list
func (s *Snips) append(message string, senderAddr string, time int) {
	mutex := &sync.Mutex{}
	mutex.Lock()

	s.snips = append(s.snips, singleSnip{message: message, senderAddr: senderAddr, timeStamp: time})
	mutex.Unlock()
}

// This function formats the list into a string and returns a list of snips received as a string
func (s *Snips) GetSnipsAsString() string {
	mutex := &sync.Mutex{}
	mutex.Lock()

	totalSnips := strconv.Itoa(len(s.snips))
	snips := totalSnips + "\n"
	for _, snip := range s.snips {
		snips += strconv.Itoa(snip.timeStamp) + " " + snip.message + " " + snip.senderAddr + "\n"
	}
	mutex.Unlock()
	return snips
}

// receive a snip and handle it
func (s *Snips) ListenForSnips(sourceAddress string, connection *net.UDPConn, ctx context.Context, peerInfo *PeerInfo) {
	listenChannel := make(chan string, 5)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			log.WithField("message", scanner.Text()).Debug("Received snip")
			listenChannel <- scanner.Text()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			log.Debug("ListenForSnips: Done")
			return
		case snip := <-listenChannel:
			sendSnipToPeers(snip, sourceAddress, connection, peerInfo)
		}
	}
}

// Send a snip to all peers
func sendSnipToPeers(content, sourceAddress string, connection *net.UDPConn, peerlist *PeerInfo) {
	log.Debug("Inside sending snip to peers")
	currentTimeStamp++
	snipTimeStamp := strconv.Itoa(currentTimeStamp)
	totalpeersSent := 0
	snipMessage := "snip" + snipTimeStamp + " " + content
	for _, peer := range peerlist.MainPeerInfo {
		if peer.IsAlive {
			protocols.SendMessage(connection, peer.PeerAddr, snipMessage)
			log.WithField("peer", peer.PeerAddr).Debug("Sent snip to peer")
			totalpeersSent++
		}

	}
	log.WithField("snipMessage", snipMessage).WithField("totalpeersSent", totalpeersSent).Debug("Sent snip to peers")
}

// Stores a new nip into the the snips list
func (s *Snips) SaveSnip(command string, senderAddr string, peerInfo *PeerInfo) {
	message := strings.Split(command, " ")

	snipTimeStamp, err := strconv.Atoi(message[0])
	if err != nil {
		log.WithField("timeStamp", message[0]).Error("TimeStamp is not a number")
		return
	}
	if len(message) < 2 {
		log.WithField("message", message[1]).WithField("command", command).Error("Invalid snip command")
		return
	}

	snipContent := strings.Join(message[1:], " ")

	// check which timestamp is bigger
	if snipTimeStamp != currentTimeStamp {
		log.WithField("snipTimeStamp", snipTimeStamp).WithField("globalTimestamp", currentTimeStamp).Debug("TimeStamp is not the same")
		currentTimeStamp = getMax(snipTimeStamp, currentTimeStamp)

	}
	// update last seen for main info
	peerInfo.UpdateMainInfoLastSeen(senderAddr)

	s.append(snipContent, senderAddr, currentTimeStamp)

}

func getMax(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
