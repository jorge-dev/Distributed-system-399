package peercommunicator

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/jorge-dev/Distributed-system-559/src/protocols"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
	log "github.com/sirupsen/logrus"
)

type PeerCommunicator struct {
	udpAddr    string
	connection *net.UDPConn
}

func NewPeerCommunicator(udpAddr string) PeerCommunicator {
	return PeerCommunicator{udpAddr: udpAddr}
}

const (
	// Define the Registry's request types

	UDP_STOP string = "stop"
	UDP_SNIP string = "snip"
	UDP_PEER string = "peer"
)

var wg sync.WaitGroup

func (pc *PeerCommunicator) Start(peerInfo *sysTypes.PeerInfo, snips *sysTypes.Snips, ctx context.Context) error {

	udpProto := protocols.NewDUP(pc.udpAddr)
	var err error
	peerInfo.AppendMainInfoPeer(pc.udpAddr, pc.udpAddr)
	pc.connection, err = udpProto.ConnectToClient()
	if err != nil {
		log.Debugf("Error while trying to connect to %s due to following error: \n %v", pc.udpAddr, err)
		return err

	}
	defer pc.connection.Close()
	PeerCommunicatorCtx, cancel := context.WithCancel(ctx)
	wg.Add(2)
	go func() {
		defer wg.Done()
		pc.messageHandler(peerInfo, snips, PeerCommunicatorCtx, cancel)
	}()
	// go func() {
	// 	defer wg.Done()
	// 	peerInfo.HandleDeadPeers(PeerCommunicatorCtx)
	// }()
	go func() {
		defer wg.Done()
		log.Debug("Inside Listening Peers")
		snips.ListenForSnips(pc.udpAddr, pc.connection, PeerCommunicatorCtx, peerInfo)
	}()
	go func() {
		defer wg.Done()
		peerInfo.SendHeartbeat(pc.connection, PeerCommunicatorCtx)
	}()

	wg.Wait()

	return nil
}

func (pc *PeerCommunicator) messageHandler(peerList *sysTypes.PeerInfo, snips *sysTypes.Snips, ctx context.Context, cancel context.CancelFunc) {
	peerMainInfo := peerList.GetMainInfoPeerList()
	// handle gracefully close
	go func() {
		<-ctx.Done()
		pc.connection.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Debug("Closing connection from messageHandler")
			return
		default:
			command, senderAdr, err := protocols.ReceiveMessage(pc.connection)
			if err != nil {
				log.Errorf("Error while trying to receive message due to following error: \n %v", err)
				return
			}
			log.WithField("command", command).WithField("senderAddress", senderAdr).Debug("Received command")
			// update last seen
			isInPeers := false
			if len(peerMainInfo) > 0 {
				for _, peer := range peerMainInfo {
					if peer.PeerAddr == senderAdr {
						peer.LastSeen = time.Now()
						peer.IsAlive = true
						isInPeers = true
						break

					}
				}
				if !isInPeers {
					peerList.AppendMainInfoPeer(senderAdr, pc.udpAddr)
				}
			}

			log.WithField("command", command).WithField("senderAddress", senderAdr).Debug("Message received")
			if len(command) > 4 {
				switch command[:4] {
				case UDP_STOP:
					log.Debug("Received stop command")
					pc.connection.Close()
					cancel()
					return
				case UDP_SNIP:
					log.Debug("Received snip command")
					command := strings.Trim(command[4:], "\n")
					go snips.SaveSnip(command, senderAdr, peerList)
				case UDP_PEER:
					log.Debug("Received peer command")
					peerAddr := strings.Trim(command[4:], "\n")
					go peerList.AppendReceivedInfoPeer(peerAddr, senderAdr)
				default:
					log.WithField("command", command).WithField("senderAddress", senderAdr).Debug("Unknown command")
				}
			} else {
				log.Error("Message is not long enough to be a valid command")

			}

		}
	}
}
