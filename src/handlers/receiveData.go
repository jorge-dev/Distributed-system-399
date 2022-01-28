package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

func ReceivePeers(scanner *bufio.Scanner) sysTypes.Peer {
	var peer sysTypes.Peer = sysTypes.NewPeer(nil, 0)
	fmt.Println("Server is sending a list of peers")
	fmt.Println("Receive it after the Enter key is pressed: ")
	bufio.NewReader(os.Stdin).ReadString('\n')
	// sources.timeStamp = getCurrentDateTime()
	// get the number of peers
	scanner.Scan()
	num, _ := strconv.Atoi(scanner.Text())
	peer.NumPeers = num
	// get the peers
	for i := 0; i < num; i++ {
		scanner.Scan()
		// insert if unique
		if strings.Contains(strings.Join(peer.GetPeerList(), " "), scanner.Text()) == false {
			peer.Append(scanner.Text())
			// peers = append(peers.p, scanner.Text())
		}
	}
	// sources.peers.peerList = peers.peerList
	fmt.Printf("Peers Received: %v\n\n", peer.GetPeerList())
	return peer
}
