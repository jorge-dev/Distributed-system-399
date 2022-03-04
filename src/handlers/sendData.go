// Handle all send requests

package handlers

import (
	"fmt"
	"net"
	"strconv"

	"github.com/jorge-dev/Distributed-system-559/src/common"
	"github.com/jorge-dev/Distributed-system-559/src/peerProc"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

func SendTeamName(conn net.Conn, teamName string) {
	fmt.Println("Server is asking for your team name ")
	// fmt.Println("Send it after the Enter key is pressed: ")
	// bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Fprintf(conn, teamName+"\n")
}

func SendCode(conn net.Conn, counter int) {
	didPrint := false
	fmt.Println("Server is, asking for your code ")
	// fmt.Println("Send it after the Enter key is pressed: ")
	// bufio.NewReader(os.Stdin).ReadString('\n')
	if counter > 0 {
		didPrint = true
	}
	code := common.PrintAllFiles("../src/", didPrint)
	fmt.Fprintf(conn, "Go\n%s\n...\n", code)
}

func SendLocation(conn net.Conn, location string) {
	fmt.Println("Server is asking for your location ")
	// fmt.Println("Send it after the Enter key is pressed: ")
	// bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Fprintf(conn, location+"\n")
}

func SendReport(conn net.Conn, peers sysTypes.Peer, sources []sysTypes.Source) {
	fmt.Println("Server is asking for your report ")
	fmt.Printf("peers: %v\n", peers.GetPeerList())
	// fmt.Println("Send it after the Enter key is pressed: ")
	// bufio.NewReader(os.Stdin).ReadString('\n')
	report := strconv.Itoa(peers.NumPeers) + "\n"
	if peers.NumPeers == 0 {
		fmt.Println("No peers found")
		report += "0\n0\n"
		conn.Write([]byte(report))

	} else {
		fmt.Println("After if inseide else")
		numberOfSources := len(sources)

		for _, peer := range peers.GetPeerList() {
			report += peer + "\n"
		}
		fmt.Println("After for inseide else")
		report += strconv.Itoa(numberOfSources) + "\n"
		report += sources[0].GetSourceAddress() + "\n"
		report += sources[0].GetTimeStamp() + "\n"

		report += peerProc.PreparelistPeersToString()
		report += peerProc.PreparelistReceivedPeerInfoToString()
		report += peerProc.PreparelistSentPeerInfoToString()
		report += peerProc.PreparelistSnipsToString()
		// fmt.Printf("\nReport: %s\n", report)
		conn.Write([]byte(report))

	}
}
