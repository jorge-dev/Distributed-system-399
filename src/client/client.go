// TCP client code

package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"

	"github.com/jorge-dev/Distributed-system-559/src/handlers"
	"github.com/jorge-dev/Distributed-system-559/src/peerProc"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

const (
	// Define the Registry's request types
	GET_NAME      string = "get team name"
	GET_CODE      string = "get code"
	GET_REPORT    string = "get report"
	GET_LOCATION  string = "get location"
	CLOSE         string = "close"
	RECEIVE_PEERS string = "receive peers"
	UDP_STOP      string = "stop"
	UDP_SNIP      string = "snip"
	UDP_PEER      string = "peer"
)

var peer sysTypes.Peer

// Creates a new client and attempts to connect to the server
func ConnectTCP(host, port, udpHost, udpPort string, ctx context.Context) error {

	//Save the host and port as a full address and initialize variables
	sourceAddress := host + ":" + port
	udpSourceAddress := udpHost + ":" + udpPort

	sources := []sysTypes.Source{sysTypes.NewSource(sourceAddress, &peer)}

	// Attempt to connect to the socket
	connection, err := net.Dial("tcp", sourceAddress)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", sourceAddress, err)
		return err
	}

	// close the connection when the function returns
	defer connection.Close()

	scanner := bufio.NewScanner(connection)
	getCodeRequestCounter := 0
	// loop until the connection is closed or error is found
	go func() {
		<-ctx.Done()
		connection.Close()
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Closing the TCP connection")
			connection.Close()
			return nil
		default:

			if !scanner.Scan() {
				fmt.Println("Server Disconnected")
				break loop
			}

			switch scanner.Text() {
			case GET_NAME:
				handlers.SendTeamName(connection, "Jorge Avila")

			case GET_CODE:
				handlers.SendCode(connection, getCodeRequestCounter)
				getCodeRequestCounter++

			case GET_LOCATION:
				handlers.SendLocation(connection, udpSourceAddress)
			case GET_REPORT:
				handlers.SendReport(connection, peer, sources)

			case RECEIVE_PEERS:
				peer = handlers.ReceivePeers(scanner, &sources[0])

			case CLOSE:
				fmt.Println("Server is closing the connection ...")
				connection.Close()
				return nil
			default:
				fmt.Printf("Unknown request &s\n", scanner.Text())

			}
		}
	}
	return nil
}

func ConnectUdpServer(host string, port string, ctx context.Context) error {

	//Save the host and port as a full address and initialize variables
	sourceAddress := host + ":" + port

	fmt.Printf("Initializing UDP server on %s\n", sourceAddress)
	udpAddr, err := net.ResolveUDPAddr("udp", sourceAddress)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", sourceAddress, err)
		return err
	}
	// checkError(err)
	conn, err := net.ListenUDP("udp", udpAddr)
	// checkError(err)
	if err != nil {
		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", sourceAddress, err)
		return err
	}

	// Call the peer process
	peerProc.PeerProcess(conn, sourceAddress, ctx)
	return err

}

// type PeerInfo struct {
// 	peerAddress   string
// 	sourceAddress string
// 	lastSeen      time.Time
// }

// type ReceivedPeerinfo struct {
// 	peerAddrReceived string
// 	peerAddrSender   string
// 	timestamp        time.Time
// }

// type SentPeerInfo struct {
// 	peerAddr     string
// 	receiverAddr string
// 	timestamp    time.Time
// }

// type Snip struct {
// 	message    string
// 	senderAddr string
// 	timeStamp  int
// }

// var listPeers []PeerInfo
// var listSnips []Snip
// var listReceivedPeerinfo []ReceivedPeerinfo
// var listSentPeerInfo []SentPeerInfo

// var mutex = &sync.Mutex{}
// var currentTime int = 0

// func AddPeer(peerAddress string, sourceAddress string) {
// 	mutex.Lock()
// 	listPeers = append(listPeers, PeerInfo{peerAddress, sourceAddress, time.Now()})
// 	mutex.Unlock()
// }

// func peerProcess(conn *net.UDPConn, sourceAddress string) {
// 	listPeers = append(listPeers, PeerInfo{sourceAddress, sourceAddress, time.Now()})
// 	fmt.Printf("Peer Party Started at %s\n", sourceAddress)
// 	wg := sync.WaitGroup{}
// 	wg.Add(4)
// 	go func() {
// 		defer wg.Done()
// 		messageHandler(conn, sourceAddress)
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		snipHandler(sourceAddress)
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		peerSender(sourceAddress)
// 	}()

// 	// go func() {
// 	// 	defer wg.Done()
// 	// 	handleInactivePeers(sourceAddress)
// 	// }()
// 	wg.Wait()

// }

// func handleInactivePeers(sourceAddress string) {
// 	for {
// 		time.Sleep(time.Second * 10)
// 		// make a copy of the list of peers
// 		listPeersCopy := make([]PeerInfo, len(listPeers))
// 		copy(listPeersCopy, listPeers)
// 		mutex.Lock()
// 		if len(listPeers) > 0 {
// 			for i := 0; i < len(listPeersCopy); i++ {
// 				if listPeersCopy[i].peerAddress != sourceAddress {
// 					if time.Since(listPeers[i].lastSeen) > time.Second*10 {
// 						listPeers = append(listPeers[:i], listPeers[i+1:]...)
// 						fmt.Printf("Peer %s is inactive and removed from the list\n", listPeers[i].peerAddress)
// 					}
// 				}
// 			}
// 			fmt.Printf("Inactive peers removed: %d\n", len(listPeers))
// 		}
// 		mutex.Unlock()
// 	}
// }

// func peerSender(sourceAddress string) {
// 	for {
// 		time.Sleep(time.Second * 5)
// 		mutex.Lock()
// 		if len(listPeers) > 0 {
// 			currentTime++
// 			for i := 0; i < len(listPeers); i++ {
// 				for j := 0; j < len(listPeers); j++ {
// 					if CheckForValidAddress(listPeers[j].peerAddress) {
// 						if listPeers[j].peerAddress != sourceAddress {

// 							msg := UDP_PEER + listPeers[i].peerAddress
// 							sendMessage(listPeers[j].peerAddress, msg)

// 							// Update Sent peer info
// 							listSentPeerInfo = append(listSentPeerInfo, SentPeerInfo{listPeers[i].peerAddress, listPeers[j].peerAddress, time.Now()})
// 						}
// 					}
// 				}
// 			}
// 		}
// 		mutex.Unlock()
// 	}

// }

// func snipHandler(sourceAddress string) {
// 	ch := make(chan string)
// 	go func() {
// 		scanner := bufio.NewScanner(os.Stdin)
// 		for scanner.Scan() {
// 			ch <- scanner.Text()
// 		}
// 	}()
// 	for {
// 		select {
// 		case msg := <-ch:
// 			sendSnip(msg, sourceAddress)
// 		}
// 	}
// }

// func sendSnip(msg string, sourceAddress string) {
// 	snipCurrentTime := strconv.Itoa(currentTime)
// 	msg = "snip" + snipCurrentTime + " " + msg
// 	currentTime++
// 	mutex.Lock()
// 	// Send the message to all peers
// 	for _, peer := range listPeers {
// 		if CheckForValidAddress(peer.peerAddress) && peer.peerAddress != sourceAddress {
// 			go sendMessage(peer.peerAddress, msg)
// 		}
// 	}
// 	mutex.Unlock()
// }

// func sendMessage(peerAddress, msg string) {
// 	conn := startUdpClient(peerAddress)
// 	defer conn.Close()
// 	fmt.Printf("Sending message to %s\n", peerAddress)
// 	_, err := conn.Write([]byte(msg))
// 	if err != nil {
// 		fmt.Printf("Error while sending message to %s due to following error: \n %v", peerAddress, err)
// 		return
// 	}

// }

// func startUdpClient(address string) net.Conn {
// 	udpAdd, err := net.ResolveUDPAddr("udp", address)
// 	if err != nil {
// 		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", address, err)
// 		return nil
// 	}
// 	conn, err := net.DialUDP("udp", nil, udpAdd)
// 	if err != nil {
// 		log.Fatalf("Error while trying to connect to %s due to following error: \n %v", address, err)
// 		return nil
// 	}
// 	return conn
// }

// func CheckForValidAddress(address string) bool {
// 	// check if the host and port are valid
// 	_, err := net.ResolveUDPAddr("udp", address)
// 	if err != nil {
// 		return false
// 	}
// 	return true
// }

// func messageHandler(conn *net.UDPConn, sourceAddress string) {
// 	for {
// 		msg, senderAddr, err := receiveUdpMessage(sourceAddress, conn)
// 		if err != nil {
// 			fmt.Println("Error while receiving message: ", err)
// 			continue
// 		}
// 		// update last seen
// 		for i := 0; i < len(listPeers); i++ {
// 			if listPeers[i].peerAddress == senderAddr {
// 				listPeers[i].lastSeen = time.Now()
// 			}
// 		}

// 		// only focus on first 4 characters
// 		fmt.Printf("Message received from %s: %s\n", senderAddr, msg)
// 		switch msg[:4] {
// 		case UDP_STOP:
// 			fmt.Println("Stopping UDP server")
// 			conn.Close()
// 			return
// 		case UDP_SNIP:
// 			fmt.Println("Snipping UDP server")
// 			command := strings.Trim(msg[4:], "\n")
// 			go storeSnips(command, senderAddr)
// 		case UDP_PEER:
// 			fmt.Println("Peer info received")
// 			peerAddr := strings.Trim(msg[4:], "\n")
// 			go storePeers(peerAddr, senderAddr)
// 		}

// 	}
// }

// func storePeers(peerAddr string, senderAddr string) {
// 	// Get the peer and source index
// 	peerIndex := peerListIndexLookUp(peerAddr)
// 	sourceIndex := peerListIndexLookUp(senderAddr)

// 	// If the peer is not in the list, add it
// 	if peerIndex == -1 {
// 		listPeers = append(listPeers, PeerInfo{peerAddr, senderAddr, time.Now()})
// 		fmt.Printf("New peer added: %s\n", peerAddr)
// 	}

// 	// If the source is not in the list, add it
// 	if sourceIndex == -1 {
// 		listPeers = append(listPeers, PeerInfo{senderAddr, senderAddr, time.Now()})
// 		fmt.Printf("New source added: %s\n", senderAddr)
// 	}

// 	listReceivedPeerinfo = append(listReceivedPeerinfo, ReceivedPeerinfo{peerAddr, senderAddr, time.Now()})

// }

// func peerListIndexLookUp(peerAddr string) int {
// 	for i := 0; i < len(listPeers); i++ {
// 		if listPeers[i].peerAddress == peerAddr {
// 			return i
// 		}

// 	}
// 	return -1

// }

// func storeSnips(command string, senderAddr string) {
// 	msg := strings.Split(command, " ")
// 	timestamp, _ := strconv.Atoi(msg[0])
// 	if len(msg) != 2 {
// 		fmt.Println("Invalid snip command")
// 		return
// 	}
// 	// Store the snip to list
// 	listSnips = append(listSnips, Snip{msg[1], senderAddr, timestamp})

// 	// update last seen
// 	for i := 0; i < len(listPeers); i++ {
// 		if listPeers[i].peerAddress == senderAddr {
// 			listPeers[i].lastSeen = time.Now()
// 		}
// 	}

// 	// check which time is the latest
// 	if timestamp > currentTime {
// 		currentTime = timestamp
// 	}

// }

// func receiveUdpMessage(address string, conn *net.UDPConn) (string, string, error) {

// 	// Read from the connection
// 	data := make([]byte, 1024)
// 	len, addr, err := conn.ReadFromUDP(data)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	msg := strings.TrimSpace(string(data[:len]))

// 	return msg, addr.String(), nil

// }
