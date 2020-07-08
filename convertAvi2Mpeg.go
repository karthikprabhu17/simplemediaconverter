package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Counter
var count uint = 1

// STATUS Type for file progress
type STATUS uint

// Status for file process progress
const (
	YETTOSTART STATUS = 0
	INPROGRESS STATUS = 1
	DONE       STATUS = 3
	FAILED     STATUS = 4
)

// AviFiles is a data structure used for storing information
type AviFiles struct {
	inFilepath  string
	outFilename string
	outputDir   string
	no          uint
	status      STATUS
	mtx         sync.Mutex
}

// ProcessingQueue is a Global var for storing all structures of files to be processed
var ProcessingQueue = []*AviFiles{}

func exit(message string, code int) {
	fmt.Println(message)
	os.Exit(code)
}

func (queueItem *AviFiles) getInputFile(state STATUS) string {
	return queueItem.inFilepath
}

func (queueItem *AviFiles) setStatus(state STATUS) {
	queueItem.status = state
}

func getOutFilename(inFilePath string) string {
	ext := filepath.Ext(inFilePath)
	outfname := inFilePath[0 : len(inFilePath)-len(ext)]

	if len(outfname) > 0 {
		return outfname
	}

	outfname = outfname + ".mp4"
	return outfname
}

func getOutputDir(inFilePath string) string {
	outputDir := filepath.Base(inFilePath)

	if len(outputDir) > 0 {
		return outputDir
	}

	return outputDir
}

func mediawalk(path string, info os.FileInfo, err error) error {
	if !info.IsDir() && strings.HasSuffix(info.Name(), "avi") {
		fmt.Printf(" %s\n", path)
	}

	aviobj := &AviFiles{
		inFilepath:  path,
		outFilename: getOutFilename(path),
		outputDir:   getOutputDir(path),
		no:          count,
		status:      YETTOSTART,
	}

	count = count + 1

	ProcessingQueue = append(ProcessingQueue, aviobj)

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
