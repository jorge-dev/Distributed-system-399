package peerProc

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
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

// convert listPeers into a string
// func ConvertlistPeersToString() string {
// 	var numPeers string = strconv.Itoa(len(listPeers))
// 	var peerList string = numPeers + "\n"
// 	for _, peer := range listPeers {
// 		peerList += peer.peerAddress + "\n"
// 	}
// 	return peerList

// }

// func ConvertlistReceivedPeerInfoToString() string {
// 	var numPeersRec string = strconv.Itoa(len(listReceivedPeerinfo))
// 	var peerListRec string = numPeersRec + "\n"
// 	for _, peer := range listReceivedPeerinfo {
// 		peerListRec += peer.peerAddrSender + " " + peer.peerAddrReceived + " " + peer.timestamp.Format("2006-01-02 15:04:05") + "\n"
// 	}
// 	return peerListRec
// }

// func ConvertlistSnipsToString() string {
// 	var numSnips string = strconv.Itoa(len(listSnips))
// 	var snipList string = numSnips + "\n"

// 	for _, snip := range listSnips {
// 		snipList += strconv.Itoa(snip.timeStamp) + " " + snip.message + " " + snip.senderAddr + "\n"
// 	}
// 	return snipList
// }

// func ConvertlistSentPeerInfoToString() string {
// 	var numPeersSent string = strconv.Itoa(len(listSentPeerInfo))
// 	var peerListSent string = numPeersSent + "\n"
// 	for _, peer := range listSentPeerInfo {
// 		peerListSent += peer.receiverAddr + " " + peer.peerAddr + " " + peer.timestamp.Format("2006-01-02 15:04:05") + "\n"
// 	}
// 	return peerListSent
// }

func AddPeer(peerAddress string, sourceAddress string) {
	// check if the peer is already in the list
	mutex.Lock()
	for _, peer := range listPeers {
		if peer.peerAddress == peerAddress {
			mutex.Unlock()
			return
		} else if peer.sourceAddress == peerAddress {
			mutex.Unlock()
			return
		} else if peer.sourceAddress == sourceAddress {
			mutex.Unlock()
			return

		}
	}
	listPeers = append(listPeers, PeerInfo{peerAddress, sourceAddress, time.Now()})
	mutex.Unlock()
	// fmt.Printf("Peers in the list: %v\n", listPeers)

	// mutex.Lock()
	// listPeers = append(listPeers, PeerInfo{peerAddress, sourceAddress, time.Now()})
	// fmt.Printf("Peers in the list: %v\n", listPeers)
	// mutex.Unlock()
}

func PeerProcess(conn *net.UDPConn, sourceAddress string, ctx context.Context) {
	listPeers = append(listPeers, PeerInfo{sourceAddress, sourceAddress, time.Now()})
	fmt.Printf("Peer Party Started at %s\n", sourceAddress)
	wg := sync.WaitGroup{}
	childCtx, cancel := context.WithCancel(ctx)
	wg.Add(4)
	go func() {
		defer wg.Done()
		messageHandler(conn, sourceAddress, childCtx, cancel)
	}()

	go func() {
		defer wg.Done()
		snipHandler(sourceAddress, conn, childCtx)
	}()

	go func() {
		defer wg.Done()
		peerSender(sourceAddress, conn, childCtx)
	}()

	go func() {
		defer wg.Done()
		handleInactivePeers(sourceAddress, childCtx)
	}()
	wg.Wait()

}

func handleInactivePeers(sourceAddress string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 15):
		}
		// time.Sleep(time.Second * 10)
		// make a copy of the list of peers
		// listPeersCopy := make([]PeerInfo, len(listPeers))
		// copy(listPeersCopy, listPeers)
		mutex.Lock()
		// fmt.Printf("Peers in the list: %v\n", listPeers)
		if len(listPeers) > 0 {
			for i := 0; i < len(listPeers); i++ {
				if listPeers[i].peerAddress != sourceAddress {
					if time.Since(listPeers[i].lastSeen) > time.Second*10 {
						listPeers = append(listPeers[:i], listPeers[i+1:]...)
						// fmt.Printf("Peer %s is inactive and removed from the list\n", listPeers[i].peerAddress)
					}
				}
			}
			fmt.Printf("Inactive peers removed. Peers left %d\n", len(listPeers))
		}
		mutex.Unlock()
	}
}

func peerSender(sourceAddress string, conn *net.UDPConn, context context.Context) {
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-context.Done():
			return
		case <-time.After(time.Second * 5):

		}

		mutex.Lock()
		if len(listPeers) > 0 {
			peerCount := 0
			// currentTime++
			// send a random peer to all peers
			peerlen := len(listPeers)
			randPeer := listPeers[rand.Intn(peerlen)]
			// fmt.Println("Sending peers")
			for _, peer := range listPeers {
				if CheckForValidAddress(peer.peerAddress) {
					sendMessage(peer.peerAddress, UDP_PEER+randPeer.peerAddress, conn)
					listSentPeerInfo = append(listSentPeerInfo, SentPeerInfo{peer.peerAddress, peer.peerAddress, time.Now()})
					peerCount++
				}
			}
			// fmt.Printf("Number of Peers sent: %d\n", peerCount)
			// for i := 0; i < len(listPeers); i++ {
			// 	if CheckForValidAddress(listPeers[j].peerAddress) {
			// 		if listPeers[i].peerAddress != sourceAddress {

			// 			msg := "peer" + listPeers[i].peerAddress + "\n"
			// 			sendMessage(listPeers[i].peerAddress, msg)

			// 			// Update Sent peer info
			// 			listSentPeerInfo = append(listSentPeerInfo, SentPeerInfo{listPeers[i].peerAddress, listPeers[j].peerAddress, time.Now()})
			// 		}

			// 	}
			// }
		}
		mutex.Unlock()
	}

}

func snipHandler(sourceAddress string, conn *net.UDPConn, ctx context.Context) {
	ch := make(chan string)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			sendSnip(msg, sourceAddress, conn)
		}
	}
}

func sendSnip(msg string, sourceAddress string, conn *net.UDPConn) {

	currentTime++
	snipCurrentTime := strconv.Itoa(currentTime)
	msg = "snip" + snipCurrentTime + " " + msg
	mutex.Lock()
	// Send the message to all peers
	// fmt.Println("Sending messages")
	for _, peer := range listPeers {
		if CheckForValidAddress(peer.peerAddress) {
			go sendMessage(peer.peerAddress, msg, conn)
		} else {
			// fmt.Printf("Invalid peer address %s\n", peer.peerAddress)
		}
	}
	mutex.Unlock()
}

func sendMessage(peerAddress, msg string, conn *net.UDPConn) {
	// startUdpClient(peerAddress,conn)
	// defer conn.Close()

	udpAdd, err := net.ResolveUDPAddr("udp", peerAddress)
	if err != nil {
		fmt.Println("Error in resolving UDP address, error is: ", err)
		return
	}

	_, err = conn.WriteToUDP([]byte(msg), udpAdd)
	// fmt.Printf("Message :like i have no one  {%s} sent to %s\n", msg, peerAddress)
	if err != nil {
		fmt.Printf("Error while sending message to %s due to following error: \n %v", peerAddress, err)
		return
	}

}

// func startUdpClient(address string conn *net.UDPConn) {
// 	udpAdd, err := net.ResolveUDPAddr("udp", address)
// 	if err != nil {
// 		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", address, err)
// 		return nil
// 	}
// 	if err != nil {
// 		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", address, err)
// 		return nil
// 	}
// 	return conn
// }

func CheckForValidAddress(address string) bool {
	// check if the host and port are valid
	_, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return false
	}
	return true
}

func messageHandler(conn *net.UDPConn, sourceAddress string, ctx context.Context, cancel context.CancelFunc) {

	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Closing the connection")
			return
		default:
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
			if len(msg) >= 4 {
				switch msg[:4] {
				case UDP_STOP:
					fmt.Println("Stopping UDP server")
					conn.Close()
					cancel()
					return
				case UDP_SNIP:
					fmt.Println("Receiving Snips")
					command := strings.Trim(msg[4:], "\n")
					go storeSnips(command, senderAddr)
				case UDP_PEER:
					fmt.Println("Peer info received")
					peerAddr := strings.Trim(msg[4:], "\n")
					go storePeers(peerAddr, senderAddr)
				default:
					fmt.Printf("Unknown command received from %s: %s\n", senderAddr, msg)

				}
			} else {
				fmt.Println("Message is not long enough to be a command")
			}
		}
	}
}

func storePeers(peerAddr string, senderAddr string) {
	// Get the peer and source index
	peerIndex := peerListIndexLookUp(peerAddr)
	sourceIndex := peerListIndexLookUp(senderAddr)

	// If the peer is not in the list, add it
	if peerIndex == -1 && sourceIndex == -1 && CheckForValidAddress(peerAddr) {
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{peerAddr, senderAddr, time.Now()})
		mutex.Unlock()
		fmt.Printf("New peer added: %s\n", peerAddr)
	} else if peerIndex == -1 && sourceIndex != -1 && CheckForValidAddress(peerAddr) {
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{peerAddr, peerAddr, time.Now()})
		mutex.Unlock()
		fmt.Printf("New peer added: %s\n", peerAddr)
	} else if peerIndex != -1 && sourceIndex == -1 && CheckForValidAddress(peerAddr) {
		mutex.Lock()
		listPeers = append(listPeers, PeerInfo{senderAddr, senderAddr, time.Now()})
	}

	// // If the source is not in the list, add it
	// if sourceIndex == -1 && CheckForValidAddress(peerAddr) {
	// 	mutex.Lock()
	// 	listPeers = append(listPeers, PeerInfo{senderAddr, senderAddr, time.Now()})
	// 	mutex.Unlock()
	// 	fmt.Printf("New source added: %s\n", senderAddr)
	// }

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
	timestamp, err := strconv.Atoi(msg[0])
	if err != nil {
		fmt.Println("Timestamp is not a valid number")
		return
	}
	if len(msg) < 2 {
		fmt.Printf("Invalid snip command: \n message: %s%s\n", command, msg)
		return
	}
	// Store the snip to list
	// join the rest of the message
	snipContent := strings.Join(msg[1:], " ")

	// check which time is the latest
	if timestamp > currentTime {
		currentTime = timestamp
	}

	mutex.Lock()
	listSnips = append(listSnips, Snip{snipContent, senderAddr, timestamp})
	mutex.Unlock()

	// update last seen
	for i := 0; i < len(listPeers); i++ {
		if listPeers[i].peerAddress == senderAddr {
			listPeers[i].lastSeen = time.Now()
		}
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
