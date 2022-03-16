// Main application entry point
package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/jorge-dev/Distributed-system-559/src/client"
	log "github.com/sirupsen/logrus"
)

func getLogLevel() log.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warn":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(getLogLevel())
	fmt.Println("Log level: ", log.GetLevel())

	log.Info("Logging initialized")

}

func main() {
	// get host and port from command line
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Usage: go run main.go <host> <port>")
		os.Exit(1)

	}

	host := args[0]
	port := args[1]

	// connect to the server
	err := client.Connect(host, port)
	if err != nil {
		fmt.Println("Error: ", err)
		log.WithField("error", err).Error("Error connecting to server")
		os.Exit(1)
	}
}
