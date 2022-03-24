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
)

var peer sysTypes.Peer
var sources []sysTypes.Source

// Creates a new client and attempts to connect to the server
func ConnectTCP(host, port, udpHost, udpPort, name string, ctx context.Context) error {

	//Save the host and port as a full address and initialize variables
	sourceAddress := host + ":" + port
	udpSourceAddress := udpHost + ":" + udpPort

	sources = []sysTypes.Source{sysTypes.NewSource(sourceAddress, &peer)}

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

	// This is the main loop of the client
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

				handlers.SendTeamName(connection, name)

			case GET_CODE:
				handlers.SendCode(connection, getCodeRequestCounter)
				getCodeRequestCounter++

			case GET_LOCATION:
				handlers.SendLocation(connection, udpSourceAddress)
			case RECEIVE_PEERS:
				peer = handlers.ReceivePeers(scanner, &sources[0])
			case GET_REPORT:
				handlers.SendReport(connection, peer, sources)
				fmt.Println("Report sent")

			case CLOSE:
				fmt.Println("Server is closing the connection ...")
				connection.Close()
				return nil
			default:
				fmt.Printf("Unknown request %s\n", scanner.Text())

			}
		}
	}
	return nil
}
func ConnectUdpServer(teamName, host string, port string, ctx context.Context) error {

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
	peerProc.PeerProcess(conn, teamName, sourceAddress, ctx)
	return err

}
