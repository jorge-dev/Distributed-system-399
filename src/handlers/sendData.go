package handlers

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/jorge-dev/Distributed-system-559/src/common"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

func SendTeamName(conn net.Conn, teamName string) {
	fmt.Println("Server is asking for your team name ")
	fmt.Println("Send it after the Enter key is pressed: ")
	bufio.NewReader(os.Stdin).ReadString('\n')
	fmt.Fprintf(conn, teamName+"\n")
}

func SendCode(conn net.Conn) {
	fmt.Println("Server is asking for your code ")
	fmt.Println("Send it after the Enter key is pressed: ")
	bufio.NewReader(os.Stdin).ReadString('\n')
	code := common.PrintAllFiles("../src/")
	fmt.Fprintf(conn, "Go\n%s\n...\n", code)
	// fmt.Fprintf(conn, "java\ncode\n...\n")
}

func SendReport(conn net.Conn, peers sysTypes.Peer, sources []sysTypes.Source) {
	fmt.Println("Server is asking for your report ")
	fmt.Println("Send it after the Enter key is pressed: ")
	bufio.NewReader(os.Stdin).ReadString('\n')
	report := strconv.Itoa(peers.NumPeers) + "\n"
	if peers.NumPeers == 0 {
		report += "0\n0\n"
		fmt.Fprintf(conn, report)
	} else {
		// This is only guaranteed to work for tis iteration
		numberOfSources := len(sources)
		peersSentFromSource := sources[0].GetPeerType()

		for _, peer := range peers.GetPeerList() {
			report += peer + "\n"
		}

		// log.Println(sources.timeStamp)
		report += strconv.Itoa(numberOfSources) + "\n"
		report += sources[0].GetSourceAddress() + "\n"
		report += sources[0].GetTimeStamp() + "\n"
		report += strconv.Itoa(peersSentFromSource.NumPeers) + "\n"
		for _, peer := range peersSentFromSource.GetPeerList() {
			report += peer + "\n"
		}

		fmt.Fprintf(conn, report)

	}

	// send the number of peers and the peers

	// fmt.Fprintf(conn, "0\n0\n0\n")
}
