package peerProc

import (
	"fmt"
	"strconv"
)

type Converter struct{}

func NewConverter() Converter {
	return Converter{}
}

func ConvertlistPeersToString() string {
	for _, peer := range listPeers {
		fmt.Println(peer.peerAddress + "\n" + peer.sourceAddress + "\n" + peer.lastSeen.Format("2006-01-02 15:04:05") + "\n")
	}

	var numPeers string = strconv.Itoa(len(listPeers))
	var peerList string = numPeers + "\n"
	for _, peer := range listPeers {
		peerList += peer.peerAddress + "\n"
	}
	return peerList

}

func ConvertlistReceivedPeerInfoToString() string {
	var numPeersRec string = strconv.Itoa(len(listReceivedPeerinfo))
	var peerListRec string = numPeersRec + "\n"
	for _, peer := range listReceivedPeerinfo {
		peerListRec += peer.peerAddrSender + " " + peer.peerAddrReceived + " " + peer.timestamp.Format("2006-01-02 15:04:05") + "\n"
	}
	return peerListRec
}

func ConvertlistSnipsToString() string {
	var numSnips string = strconv.Itoa(len(listSnips))
	var snipList string = numSnips + "\n"

	for _, snip := range listSnips {
		snipList += strconv.Itoa(snip.timeStamp) + " " + snip.message + " " + snip.senderAddr + "\n"
	}
	return snipList
}

func ConvertlistSentPeerInfoToString() string {
	var numPeersSent string = strconv.Itoa(len(listSentPeerInfo))
	var peerListSent string = numPeersSent + "\n"
	for _, peer := range listSentPeerInfo {
		peerListSent += peer.receiverAddr + " " + peer.peerAddr + " " + peer.timestamp.Format("2006-01-02 15:04:05") + "\n"
	}
	return peerListSent
}
