// import packages for the client

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func main() {
	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: client host port")
		os.Exit(1)

	}
	host := args[0]
	port := args[1]

	teamName := "Jorge Avila"
	peers := []string{}
	numPeers := 0

	// connect to the socket
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatalln("Error while trying to connect to", err)
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

		if scanner.Text() == "get team name" {
			fmt.Print("Server is asking for your team name ")
			fmt.Print("Send it after the Enter key is pressed ")
			bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Fprintf(conn, teamName+"\n")
		} else if scanner.Text() == "get code" {
			fmt.Print("Server is asking for your code ")
			fmt.Print("Send it after the Enter key is pressed ")
			bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Fprintf(conn, "java\ncode\n...\n")
		} else if scanner.Text() == "receive peers" {
			fmt.Print("Server is sending a list of peers")
			// get the number of peers
			scanner.Scan()
			num, _ := strconv.Atoi(scanner.Text())
			numPeers = num
			// get the peers
			for i := 0; i < numPeers; i++ {
				scanner.Scan()
				peers = append(peers, scanner.Text())
			}
			fmt.Println("Peers: ", peers)
		} else if scanner.Text() == "get report" {
			fmt.Print("Server is asking for your report ")
			fmt.Print("Send it after the Enter key is pressed ")
			bufio.NewReader(os.Stdin).ReadString('\n')
			// send the number of peers and the peers
			fmt.Fprintf(conn, "2\n136.159.5.27:41\n136.99.21.5:567\n1\n136.159.5.27:55921\n2021-01-25 15:18:23\n2\n136.159.5.27:41\n136.99.21.5:567\n")
		} else if scanner.Text() == "close" {
			fmt.Println("Server is closing the connection")
			break
		} else {
			fmt.Println("Server: ", scanner.Text())
		}

		// // print server response
		// fmt.Println(scanner.Text())

		// // read in input from stdin
		// userInput := bufio.NewReader(os.Stdin)

		// // read in the userInput
		// text, error := userInput.ReadString('\n')
		// if error != nil {
		// 	log.Fatalln("Error while reading user input ", error)

		// }
		// // send the text to the server
		// fmt.Fprintf(conn, text)

	}

}
