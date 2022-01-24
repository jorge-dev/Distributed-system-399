// import packages for the client

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// connect to the socket
	conn, err := net.Dial("tcp", ":55921")
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

		// print server response
		fmt.Println(scanner.Text())

		// read in input from stdin
		userInput := bufio.NewReader(os.Stdin)

		// read in the userInput
		text, error := userInput.ReadString('\n')
		if error != nil {
			log.Fatalln("Error while reading user input ", error)

		}
		// send the text to the server
		fmt.Fprintf(conn, text)

	}

}
