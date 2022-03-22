// TCP client code

package registry

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/jorge-dev/Distributed-system-559/src/handlers"
	"github.com/jorge-dev/Distributed-system-559/src/protocols"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
	log "github.com/sirupsen/logrus"
)

const (
	// Define the Registry's request types
	GET_NAME      string = "get team name"
	GET_CODE      string = "get code"
	GET_REPORT    string = "get report"
	GET_LOCATION  string = "get location"
	CLOSE         string = "close"
	RECEIVE_PEERS string = "receive peers"
)

var wg sync.WaitGroup

type Client struct {
	// Define the client's data structure
	tcpAddr  string
	udpAddr  string
	teamName string
	// peer       *sysTypes.Peer
	// sources    []sysTypes.Source
	connection *net.TCPConn
}

func NewClient(tcpAddress, udpAddress, team string) *Client {
	client := Client{
		tcpAddr:  tcpAddress,
		udpAddr:  udpAddress,
		teamName: team,
	}
	return &client
}

func (c *Client) Start(ctx context.Context) error {
	tcpProto := protocols.NewTCP(c.tcpAddr)
	var err error

	// Connect to the server
	c.connection, err = tcpProto.ConnectToClient()
	// c.connection, err := tcpProto.ConnectToClient()
	if err != nil {
		log.Debugf("Error while trying to connect to %s due to following error: \n %v", c.tcpAddr, err)
		return err
	}
	defer c.connection.Close()

	log.Info("Connection to registry has started")

	// Init a client context
	clientCtx, cancel := context.WithCancel(ctx)
	c.requestManager(clientCtx, cancel)

	return nil

}

var peer sysTypes.Peer
var sources []sysTypes.Source

func (c *Client) requestManager(ctx context.Context, cancel context.CancelFunc) {
	sources = []sysTypes.Source{sysTypes.NewSource(c.tcpAddr, &peer)}
	scanner := bufio.NewScanner(c.connection)
	wg.Add(1)
	go func() {
		<-ctx.Done()
		log.Debug("Closing request manager from Go Func")
		c.connection.Close()
		wg.Done()
	}()

	getCodeRequestCounter := 0

mainLoop:
	for {
		select {
		case <-ctx.Done():
			log.Debug("Closing request manager form loop")
			fmt.Println("Closing the TCP connection")
			c.connection.Close()
			return
		default:

			if !scanner.Scan() {
				fmt.Println("Server Disconnected")
				break mainLoop
			}
			log.WithField("message", scanner.Text()).Debug("Received message from registry")
			switch scanner.Text() {
			case GET_NAME:

				handlers.SendTeamName(c.connection, c.teamName)

			case GET_CODE:
				handlers.SendCode(c.connection, getCodeRequestCounter)
				getCodeRequestCounter++

			// case GET_LOCATION:
			// 	handlers.SendLocation(c.connection, c.udpAddr)
			case RECEIVE_PEERS:
				peer = handlers.ReceivePeers(scanner, &sources[0])
			case GET_REPORT:
				handlers.SendReport(c.connection, peer, sources)
				fmt.Println("Report sent")

			case CLOSE:
				log.Debug("Closing request manager due to a close request")
				// fmt.Println("Server is closing the connection ...")
				c.connection.Close()
				return
			default:
				log.Errorf("Unknown request received: %s", scanner.Text())
				// fmt.Printf("Unknown request %s\n", scanner.Text())

			}
		}
	}

	cancel()
	wg.Wait()
	log.Warn("Request manager has been closed")

}
