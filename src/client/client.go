// import packages for the client

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Peer struct {
	address     string
	peer        []string
	numPeers    int
	source_date string
}

func getCurrentDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func listAllFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, listAllFiles(dir+file.Name()+"/")...)
		} else {
			paths = append(paths, dir+file.Name())
		}
	}
	return paths

}

func getFileContents(file_path string) string {
	file, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func printAllFiles(dir string) string {
	files := listAllFiles(dir)
	var code string

	for i, file := range files {
		code += "\n// File " + strconv.Itoa(i+1) + " out of " + strconv.Itoa(len(files)) + ":\n\n"
		code += getFileContents(file)
	}
	return code
}

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
	source := 0

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
			fmt.Println("Server is asking for your team name ")
			fmt.Println("Send it after the Enter key is pressed: ")
			bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Fprintf(conn, teamName+"\n")
		} else if scanner.Text() == "get code" {
			fmt.Println("Server is asking for your code ")
			fmt.Println("Send it after the Enter key is pressed: ")
			bufio.NewReader(os.Stdin).ReadString('\n')
			code := printAllFiles("../../src/")
			fmt.Fprintf(conn, "Go\n%s\n...\n", code)
			// fmt.Fprintf(conn, "java\ncode\n...\n")

		} else if scanner.Text() == "receive peers" {
			fmt.Println("Server is sending a list of peers")
			fmt.Println("Receive it after the Enter key is pressed: ")
			bufio.NewReader(os.Stdin).ReadString('\n')

			// get the number of peers
			scanner.Scan()
			num, _ := strconv.Atoi(scanner.Text())
			numPeers = num
			// get the peers
			for i := 0; i < numPeers; i++ {
				scanner.Scan()
				// insert if unique
				if strings.Contains(strings.Join(peers, " "), scanner.Text()) == false {
					peers = append(peers, scanner.Text())
				}
			}
			fmt.Printf("Peers Received: %v\n\n", peers)

		} else if scanner.Text() == "get report" {
			fmt.Println("Server is asking for your report ")
			fmt.Println("Send it after the Enter key is pressed: ")
			bufio.NewReader(os.Stdin).ReadString('\n')
			report := strconv.Itoa(numPeers) + "\n"
			if numPeers == 0 {
				report += "0\n0\n"
				fmt.Fprintf(conn, report)
			} else {
				// This is only guaranteed to work for tis iteration
				source = 1

				for _, peer := range peers {
					report += peer + "\n"
				}

				report += strconv.Itoa(source) + "\n"
				report += host + ":" + port + "\n"
				report += getCurrentDateTime() + "\n"
				report += strconv.Itoa(numPeers) + "\n"
				for _, peer := range peers {
					report += peer + "\n"
				}

				fmt.Fprintf(conn, report)

			}

			// send the number of peers and the peers

			// fmt.Fprintf(conn, "0\n0\n0\n")
		} else if scanner.Text() == "close" {
			fmt.Println("Server is closing the connection")
			break
		} else {
			fmt.Println("Server: ", scanner.Text())
		}

	}
	fmt.Println("Connection closed")

}
