// Helper methods used by other packages

package common

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
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
