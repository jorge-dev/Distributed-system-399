package common

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

func GetCurrentDateTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

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

func PrintAllFiles(dir string) string {
	files := ListAllFiles(dir)
	var code string

	for i, file := range files {
		code += "\n// File " + strconv.Itoa(i+1) + " out of " + strconv.Itoa(len(files)) + ":\n\n"
		code += GetFileContents(file)
	}
	return code
}
