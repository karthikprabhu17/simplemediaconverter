package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type STATUS uint

const (
	YETTOSTART = 0 STATUS
	INPROGRESS = 1
	DONE       = 3
	FAILED     = 4
)

type aviFiles struct {
	inputFilepath  string
	outputFilename string
	outputDir      string
	no             uint
	status
}

func exit(message string, code int) {
	fmt.Println(message)
	os.Exit(code)
}

func mediawalk(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.HasSuffix(info.Name(), "avi") {
		fmt.Printf(" %s\n", path)
	}

	return nil
}

func main() {
	args := os.Args

	if len(args) < 2 {
		exit("Missing input Directory arg", -1)
	}

	path := args[1]

	if _, err := os.Stat(path); os.IsNotExist(err) {
		exit("directory does not exit", 2)
	}

	fmt.Println("*************************")
	fmt.Println("Avi Files to be Processed")
	fmt.Println("*************************")
	filepath.Walk(path, mediawalk)

}
