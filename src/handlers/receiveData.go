// Handles the receiving of data from the server

package handlers

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

func ReceivePeers(scanner *bufio.Scanner, source *sysTypes.Source, peerInfo *sysTypes.PeerInfo) sysTypes.Peer {
	var peer sysTypes.Peer = sysTypes.NewPeer(nil, 0)
	source.UpdateTimeStamp()
	fmt.Println("Server is sending a list of peers")
	// bufio.NewReader(os.Stdin).ReadString('\n')

	scanner.Scan()
	num, _ := strconv.Atoi(scanner.Text())
	peer.NumPeers = num

	// get the peers
	for i := 0; i < num; i++ {
		scanner.Scan()
		// Append only if string is not already in the list
		if strings.Contains(strings.Join(peer.GetPeerList(), " "), scanner.Text()) == false {
			peer.Append(scanner.Text())
			peerInfo.AppendMainInfoPeer(scanner.Text(), source.GetAddress())

		}
	}

	fmt.Printf("\nUpdated Source TimeStamp at: %s\n\n", source.GetTimeStamp())
	fmt.Printf("Peers Received: %v\n\n", peer.GetPeerList())
	return peer
}
