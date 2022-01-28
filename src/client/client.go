// import packages for the client

package client

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/jorge-dev/Distributed-system-559/src/handlers"
	"github.com/jorge-dev/Distributed-system-559/src/sysTypes"
)

// type Peer struct {
// 	peerList []string
// 	numPeers int
// }

// func NewPeer(peerList []string, numPeers int) Peer {
// 	return Peer{peerList, numPeers}
// }

// type Source struct {
// 	address   string
// 	peers     *Peer
// 	timeStamp string
// }

// func getCurrentDateTime() string {
// 	return time.Now().Format("2006-01-02 15:04:05")
// }

// func NewSource(address string, peers *Peer) Source {
// 	return Source{address, peers, getCurrentDateTime()}
// }

// func listAllFiles(dir string) []string {
// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var paths []string
// 	for _, file := range files {
// 		if file.IsDir() {
// 			paths = append(paths, listAllFiles(dir+file.Name()+"/")...)
// 		} else {
// 			paths = append(paths, dir+file.Name())
// 		}
// 	}
// 	return paths

// }

// func getFileContents(file_path string) string {
// 	file, err := os.Open(file_path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	b, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return string(b)
// }

// func printAllFiles(dir string) string {
// 	files := listAllFiles(dir)
// 	var code string

// 	for i, file := range files {
// 		code += "\n// File " + strconv.Itoa(i+1) + " out of " + strconv.Itoa(len(files)) + ":\n\n"
// 		code += getFileContents(file)
// 	}
// 	return code
// }

const (
	// Define the Registry's request types
	GET_NAME      string = "get team name"
	GET_CODE      string = "get code"
	GET_REPORT    string = "get report"
	CLOSE         string = "close"
	RECEIVE_PEERS string = "receive peers"
)

func Connect(host string, port string) error {

	// // get host and port from command line
	// args := os.Args[1:]
	// if len(args) != 2 {
	// 	fmt.Println("Usage: client host port")
	// 	os.Exit(1)

	// }

	// host := args[0]
	// port := args[1]
	// address := host + ":" + port

	// teamName := "Jorge Avila"
	// var peers Peer
	// sources := NewSource(address, &peers)
	// fmt.Println(sources.timeStamp)
	// // peers := []string{}
	// // numPeers := 0
	// source := 0
	sourceAddress := host + ":" + port
	var peer sysTypes.Peer
	sources := []sysTypes.Source{sysTypes.NewSource(sourceAddress, &peer)}
	// connect to the socket
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalln("Error while trying to connect to", err)
		return err
	}
	// close the connection when the function returns
	defer conn.Close()

	// create a new scanner
	scanner := bufio.NewScanner(conn)
	// loop through the scanner

	for {

		if !scanner.Scan() {
			fmt.Println("Server Disconnected")
			break
		}

		switch scanner.Text() {
		case GET_NAME:
			handlers.SendTeamName(conn, "Jorge Avila")
			break
		case GET_CODE:
			handlers.SendCode(conn)
			break
		case GET_REPORT:
			handlers.SendReport(conn, peer, sources)
			break
		case RECEIVE_PEERS:
			peer = handlers.ReceivePeers(scanner)
			break
		case CLOSE:
			fmt.Println("Server is closing the connection")
			break
		default:
			fmt.Printf("Unknown request &s\n", scanner.Text())
			break
		}

		// for {
		// 	if !scanner.Scan() {
		// 		fmt.Println("Server Disconnected")
		// 		break
		// 	}

		// 	if scanner.Text() == "get team name" {
		// 		fmt.Println("Server is asking for your team name ")
		// 		fmt.Println("Send it after the Enter key is pressed: ")
		// 		bufio.NewReader(os.Stdin).ReadString('\n')
		// 		fmt.Fprintf(conn, teamName+"\n")
		// 	} else if scanner.Text() == "get code" {
		// 		fmt.Println("Server is asking for your code ")
		// 		fmt.Println("Send it after the Enter key is pressed: ")
		// 		bufio.NewReader(os.Stdin).ReadString('\n')
		// 		code := printAllFiles("../../src/")
		// 		fmt.Fprintf(conn, "Go\n%s\n...\n", code)
		// 		// fmt.Fprintf(conn, "java\ncode\n...\n")

		// 	} else if scanner.Text() == "receive peers" {
		// 		fmt.Println("Server is sending a list of peers")
		// 		fmt.Println("Receive it after the Enter key is pressed: ")
		// 		bufio.NewReader(os.Stdin).ReadString('\n')
		// 		// sources.timeStamp = getCurrentDateTime()
		// 		// get the number of peers
		// 		scanner.Scan()
		// 		num, _ := strconv.Atoi(scanner.Text())
		// 		peers.numPeers = num
		// 		// get the peers
		// 		for i := 0; i < peers.numPeers; i++ {
		// 			scanner.Scan()
		// 			// insert if unique
		// 			if strings.Contains(strings.Join(peers.peerList, " "), scanner.Text()) == false {
		// 				peers.peerList = append(peers.peerList, scanner.Text())
		// 				// peers = append(peers.p, scanner.Text())
		// 			}
		// 		}
		// 		// sources.peers.peerList = peers.peerList
		// 		fmt.Printf("Peers Received: %v\n\n", peers.peerList)

		// 	} else if scanner.Text() == "get report" {
		// 		fmt.Println("Server is asking for your report ")
		// 		fmt.Println("Send it after the Enter key is pressed: ")
		// 		bufio.NewReader(os.Stdin).ReadString('\n')
		// 		report := strconv.Itoa(peers.numPeers) + "\n"
		// 		if peers.numPeers == 0 {
		// 			report += "0\n0\n"
		// 			fmt.Fprintf(conn, report)
		// 		} else {
		// 			// This is only guaranteed to work for tis iteration
		// 			source = 1

		// 			for _, peer := range peers.peerList {
		// 				report += peer + "\n"
		// 			}
		// 			// log.Println(sources.timeStamp)
		// 			report += strconv.Itoa(source) + "\n"
		// 			report += sources.address + "\n"
		// 			report += sources.timeStamp + "\n"
		// 			report += strconv.Itoa(sources.peers.numPeers) + "\n"
		// 			for _, peer := range sources.peers.peerList {
		// 				report += peer + "\n"
		// 			}

		// 			fmt.Fprintf(conn, report)

		// 		}

		// 		// send the number of peers and the peers

		// 		// fmt.Fprintf(conn, "0\n0\n0\n")
		// 	} else if scanner.Text() == "close" {
		// 		fmt.Println("Server is closing the connection")
		// 		break
		// 	} else {
		// 		fmt.Println("Server: ", scanner.Text())
		// 	}

	}
	fmt.Println("Connection closed")

	return nil

}
