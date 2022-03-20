// TCP client code

package client

import (
	"bufio"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/jorge-dev/Distributed-system-559/src/handlers"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

const (
	// Define the Registry's request types
	GET_NAME      string = "get team name"
	GET_CODE      string = "get code"
	GET_REPORT    string = "get report"
	CLOSE         string = "close"
	RECEIVE_PEERS string = "receive peers"
)

// Creates a new client and attempts to connect to the server
func Connect(host string, port string) error {

	//Save the host and port as a full address and initialize variables
	sourceAddress := host + ":" + port
	var peer sysTypes.Peer
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
	for {

		if !scanner.Scan() {
			fmt.Println("Server Disconnected")
			break
		}

		switch scanner.Text() {
		case GET_NAME:
			handlers.SendTeamName(connection, "Jorge Avila")
			break
		case GET_CODE:
			handlers.SendCode(connection, getCodeRequestCounter)
			getCodeRequestCounter++
			break
		case GET_REPORT:
			handlers.SendReport(connection, peer, sources)
			break
		case RECEIVE_PEERS:
			peer = handlers.ReceivePeers(scanner, &sources[0])
			break
		case CLOSE:
			fmt.Println("Server is closing the connection ...")
			break
		default:
			fmt.Printf("Unknown request &s\n", scanner.Text())
			break
		}
	}

	fmt.Println("Connection closed")

	return nil

}
