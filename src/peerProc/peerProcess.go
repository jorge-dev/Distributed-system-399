package peerProc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// Define the Registry's request types

	UDP_STOP string = "stop"
	UDP_SNIP string = "snip"
	UDP_PEER string = "peer"
)

type PeerInfo struct {
	peerAddress   string
	sourceAddress string
	lastSeen      time.Time
}

type ReceivedPeerinfo struct {
	peerAddrReceived string
	peerAddrSender   string
	timestamp        time.Time
}

type SentPeerInfo struct {
	peerAddr     string
	receiverAddr string
	timestamp    time.Time
}

type Snip struct {
	message    string
	senderAddr string
	timeStamp  int
}

var listPeers []PeerInfo
var listSnips []Snip
var listReceivedPeerinfo []ReceivedPeerinfo
var listSentPeerInfo []SentPeerInfo

var mutex = &sync.Mutex{}
var currentTime int = 0

func AddPeer(peerAddress string, sourceAddress string) {
	mutex.Lock()
	listPeers = append(listPeers, PeerInfo{peerAddress, sourceAddress, time.Now()})
	mutex.Unlock()
}

func PeerProcess(conn *net.UDPConn, sourceAddress string) {
	listPeers = append(listPeers, PeerInfo{sourceAddress, sourceAddress, time.Now()})
	fmt.Printf("Peer Party Started at %s\n", sourceAddress)
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		messageHandler(conn, sourceAddress)
	}()

	go func() {
		defer wg.Done()
		snipHandler(sourceAddress)
	}()

	go func() {
		defer wg.Done()
		peerSender(sourceAddress)
	}()

	// go func() {
	// 	defer wg.Done()
	// 	handleInactivePeers(sourceAddress)
	// }()
	wg.Wait()

}

func handleInactivePeers(sourceAddress string) {
	for {
		time.Sleep(time.Second * 10)
		// make a copy of the list of peers
		listPeersCopy := make([]PeerInfo, len(listPeers))
		copy(listPeersCopy, listPeers)
		mutex.Lock()
		if len(listPeers) > 0 {
			for i := 0; i < len(listPeersCopy); i++ {
				if listPeersCopy[i].peerAddress != sourceAddress {
					if time.Since(listPeers[i].lastSeen) > time.Second*10 {
						listPeers = append(listPeers[:i], listPeers[i+1:]...)
						fmt.Printf("Peer %s is inactive and removed from the list\n", listPeers[i].peerAddress)
					}
				}
			}
			fmt.Printf("Inactive peers removed: %d\n", len(listPeers))
		}
		mutex.Unlock()
	}
}

func peerSender(sourceAddress string) {
	for {
		time.Sleep(time.Second * 5)
		mutex.Lock()
		if len(listPeers) > 0 {
			currentTime++
			for i := 0; i < len(listPeers); i++ {
				for j := 0; j < len(listPeers); j++ {
					if CheckForValidAddress(listPeers[j].peerAddress) {
						if listPeers[j].peerAddress != sourceAddress {

							msg := UDP_PEER + listPeers[i].peerAddress
							sendMessage(listPeers[j].peerAddress, msg)

							// Update Sent peer info
							listSentPeerInfo = append(listSentPeerInfo, SentPeerInfo{listPeers[i].peerAddress, listPeers[j].peerAddress, time.Now()})
						}
					}
				}
			}
		}
		mutex.Unlock()
	}

}

func snipHandler(sourceAddress string) {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	for {
		select {
		case msg := <-ch:
			sendSnip(msg, sourceAddress)
		}
	}
}

func sendSnip(msg string, sourceAddress string) {
	snipCurrentTime := strconv.Itoa(currentTime)
	msg = "snip" + snipCurrentTime + " " + msg
	currentTime++
	mutex.Lock()
	// Send the message to all peers
	for _, peer := range listPeers {
		if CheckForValidAddress(peer.peerAddress) && peer.peerAddress != sourceAddress {
			go sendMessage(peer.peerAddress, msg)
		}
	}
	mutex.Unlock()
}

func sendMessage(peerAddress, msg string) {
	conn := startUdpClient(peerAddress)
	defer conn.Close()
	fmt.Printf("Sending message to %s\n", peerAddress)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error while sending message to %s due to following error: \n %v", peerAddress, err)
		return
	}

}

func startUdpClient(address string) net.Conn {
	udpAdd, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", address, err)
		return nil
	}
	conn, err := net.DialUDP("udp", nil, udpAdd)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", address, err)
		return nil
	}
	return conn
}

func CheckForValidAddress(address string) bool {
	// check if the host and port are valid
	_, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return false
	}
	return true
}

func messageHandler(conn *net.UDPConn, sourceAddress string) {
	for {
		msg, senderAddr, err := receiveUdpMessage(sourceAddress, conn)
		if err != nil {
			fmt.Println("Error while receiving message: ", err)
			continue
		}
		// update last seen
		for i := 0; i < len(listPeers); i++ {
			if listPeers[i].peerAddress == senderAddr {
				listPeers[i].lastSeen = time.Now()
			}
		}

		// only focus on first 4 characters
		fmt.Printf("Message received from %s: %s\n", senderAddr, msg)
		switch msg[:4] {
		case UDP_STOP:
			fmt.Println("Stopping UDP server")
			conn.Close()
			return
		case UDP_SNIP:
			fmt.Println("Snipping UDP server")
			command := strings.Trim(msg[4:], "\n")
			go storeSnips(command, senderAddr)
		case UDP_PEER:
			fmt.Println("Peer info received")
			peerAddr := strings.Trim(msg[4:], "\n")
			go storePeers(peerAddr, senderAddr)
		}

	}
}

func storePeers(peerAddr string, senderAddr string) {
	// Get the peer and source index
	peerIndex := peerListIndexLookUp(peerAddr)
	sourceIndex := peerListIndexLookUp(senderAddr)

	// If the peer is not in the list, add it
	if peerIndex == -1 {
		listPeers = append(listPeers, PeerInfo{peerAddr, senderAddr, time.Now()})
		fmt.Printf("New peer added: %s\n", peerAddr)
	}

	// If the source is not in the list, add it
	if sourceIndex == -1 {
		listPeers = append(listPeers, PeerInfo{senderAddr, senderAddr, time.Now()})
		fmt.Printf("New source added: %s\n", senderAddr)
	}

	listReceivedPeerinfo = append(listReceivedPeerinfo, ReceivedPeerinfo{peerAddr, senderAddr, time.Now()})

}

func peerListIndexLookUp(peerAddr string) int {
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].peerAddress == peerAddr {
			return i
		}

	}
	return -1

}

func storeSnips(command string, senderAddr string) {
	msg := strings.Split(command, " ")
	timestamp, _ := strconv.Atoi(msg[0])
	if len(msg) != 2 {
		fmt.Println("Invalid snip command")
		return
	}
	// Store the snip to list
	listSnips = append(listSnips, Snip{msg[1], senderAddr, timestamp})

	// update last seen
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].peerAddress == senderAddr {
			listPeers[i].lastSeen = time.Now()
		}
	}

	// check which time is the latest
	if timestamp > currentTime {
		currentTime = timestamp
	}

}

func receiveUdpMessage(address string, conn *net.UDPConn) (string, string, error) {

	// Read from the connection
	data := make([]byte, 1024)
	len, addr, err := conn.ReadFromUDP(data)
	if err != nil {
		return "", "", err
	}
	msg := strings.TrimSpace(string(data[:len]))

	return msg, addr.String(), nil

}
