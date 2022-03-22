// Helper methods used by other packages

package common

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
)

// gets the current date and time in MT format
func GetCurrentDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// Saves all the files(not dir) in the directory to a string
func ListAllFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, ListAllFiles(dir+file.Name()+"/")...)
		} else {
			paths = append(paths, dir+file.Name())
		}
	}
	return paths

}

// Reads and stores file contents (from the path arg provided)into a string
func GetFileContents(file_path string) string {
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

// Prints all the file contents from a given directory
func PrintAllFiles(dir string, didPrint bool) string {
	files := ListAllFiles(dir)
	var code string

	for i, file := range files {
		if didPrint {
			code += "\n//****Duplicated code file due to server multiple requests****\n"
		}
		code += "\n// File " + strconv.Itoa(i+1) + " out of " + strconv.Itoa(len(files)) + ": " + file + "\n\n"
		code += GetFileContents(file)
	}
	return code
}

func getEnvVariable(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %v", err)
		return ""
	}
	return os.Getenv(key)
}

func GetlogLevel() log.Level {

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

func GetIpAndTeam(flag, udpPort string) (string, string, string) {
	if strings.Contains(flag, "t") && !strings.Contains(flag, "n") {
		testTcpAddr := getEnvVariable("TEST_HOST") + ":" + getEnvVariable("TEST_PORT")
		testUdpAddr := getEnvVariable("TEST_HOST") + ":" + udpPort
		teamName := getEnvVariable("TEAM_NAME")
		return testTcpAddr, testUdpAddr, teamName
	} else if strings.Contains(flag, "t") && strings.Contains(flag, "n") {
		testTcpAddr := getEnvVariable("TEST_HOST") + ":" + getEnvVariable("TEST_PORT")
		testUdpAddr := getEnvVariable("TEST_HOST") + ":" + udpPort
		teamName := getRandTeamName()
		return testTcpAddr, testUdpAddr, teamName
	} else if strings.Contains(flag, "r") && strings.Contains(flag, "n") {
		submissionTcpAddr := getEnvVariable("REGISTRY_HOST") + ":" + getEnvVariable("REGISTRY_PORT")
		submissionUdpAddr := getEnvVariable("REGISTRY_HOST") + ":" + udpPort
		teamName := getRandTeamName()
		return submissionTcpAddr, submissionUdpAddr, teamName
	}

	submissionTcpAddr := getEnvVariable("REGISTRY_HOST") + ":" + getEnvVariable("REGISTRY_PORT")
	submissionUdpAddr := getEnvVariable("REGISTRY_HOST") + ":" + udpPort
	teamName := getEnvVariable("TEAM_NAME")
	return submissionTcpAddr, submissionUdpAddr, teamName

}

func getRandTeamName() string {
	rand.Seed(time.Now().UnixNano())
	teamName := "Jorge Avila" + strconv.Itoa(rand.Intn(100))
	return teamName
}
