// Main application entry point
package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jorge-dev/Distributed-system-559/src/client"
	log "github.com/sirupsen/logrus"
)

func getEnvVariable(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %v", err)
		return ""
	}
	return os.Getenv(key)
}

func getlogLevel() log.Level {

	logLevel := getEnvVariable("LOG_LEVEL")

	switch logLevel {
	case "debug":
		fmt.Println("Log level set to debug")
		return log.DebugLevel
	case "info":
		fmt.Println("Log level set to info")
		return log.InfoLevel
	case "warn":
		fmt.Println("Log level set to warn")
		return log.WarnLevel
	case "error":
		fmt.Println("Log level set to error")
		return log.ErrorLevel
	case "fatal":
		fmt.Println("Log level set to fatal")
		return log.FatalLevel
	case "panic":
		fmt.Println("Log level set to panic")
		return log.PanicLevel
	default:
		fmt.Println("Log level set to info")
		return log.InfoLevel
	}

}

func getIPAddress(flag, udpPort string) (string, string) {
	if flag == "--test" || flag == "-t" {
		testTcpAddr := getEnvVariable("TEST_HOST") + ":" + getEnvVariable("TEST_PORT")
		testUdpAddr := getEnvVariable("TEST_HOST") + ":" + udpPort
		return testTcpAddr, testUdpAddr
	}
	submissionTcpAddr := getEnvVariable("REGISTRY_HOST") + ":" + getEnvVariable("REGISTRY_PORT")
	submissionUdpAddr := getEnvVariable("REGISTRY_HOST") + ":" + udpPort
	return submissionTcpAddr, submissionUdpAddr

}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(getlogLevel())
}

func main() {

	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <udpPort> <[-t | --test] || [--registry | -r]>")
		os.Exit(1)

	}

	// tcpHost := args[0]
	// tcpPort := args[1]

	// udpHost := args[2]
	udpPort := args[0]
	flag := args[1]

	tcpAddr, udpAddr := getIPAddress(flag, udpPort)

	log.WithField("tcpAddr", tcpAddr).WithField("udpAddr", udpAddr).Info("Connecting to the server")

	// connect to the server
	go func() {
		defer wg.Done()
		err := client.ConnectTCP(tcpAddr, udpAddr)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	}()

}
