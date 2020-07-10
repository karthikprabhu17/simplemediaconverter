package main

import (
	"fmt"
	"log"
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

var STATUS_STRING_MAP = map[STATUS]string{YETTOSTART: "NOTSTARTED", INPROGRESS: "INPROGRESS", DONE: "DONE", FAILED: "FAILED"}

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

func (queueItem *AviFiles) getInputFile() string {
	return queueItem.inFilepath
}

func (queueItem *AviFiles) getStatus() STATUS {
	return queueItem.status
}

func (queueItem *AviFiles) setStatus(state STATUS) {
	queueItem.status = state
}

func getOutFilename(inFilePath string) string {
	ext := filepath.Ext(inFilePath)
	outfname := inFilePath[0 : len(inFilePath)-len(ext)]

	if len(outfname) == 0 {
		log.Printf("outfilename could be determined\n")
	}

	outfname = outfname + ".mp4"
	return outfname
}

func getOutputDir(inFilePath string) string {
	outputDir := filepath.Dir(inFilePath)

	if len(outputDir) > 0 {
		return outputDir
	}

	return outputDir
}

func mediawalk(path string, info os.FileInfo, err error) error {
	var aviobj *AviFiles

	if !info.IsDir() && strings.HasSuffix(info.Name(), "avi") {
		fmt.Printf("%s\n", path)

		aviobj = &AviFiles{
			inFilepath:  path,
			outFilename: getOutFilename(path),
			outputDir:   getOutputDir(path),
			no:          count,
			status:      YETTOSTART,
		}

		log.Printf("aviobj: inFilepath: %s, outFilename: %s, outputDir: %s, no: %d\n\n",
			aviobj.inFilepath, aviobj.outFilename, aviobj.outputDir, aviobj.no)

		count = count + 1
	}

	if aviobj != nil {
		ProcessingQueue = append(ProcessingQueue, aviobj)
	}

	return nil
}

func printReport() {
	fmt.Println()
	for i := range ProcessingQueue {
		fmt.Printf("%s <----> [%s]\n", ProcessingQueue[i].inFilepath, STATUS_STRING_MAP[ProcessingQueue[i].status])
	}
}

func main() {
	logFile, err := os.OpenFile("gomediacomverter.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to create logfile, error: %s\n", err.Error())
	}

	defer logFile.Close()

	log.SetOutput(logFile)
	log.Println("log file works")

	args := os.Args

	if len(args) < 2 {
		log.Output(1, "Missing input Directory arg")
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

	for i := 0; i < 10; i++ {
		ProcessingQueue[i].runConversion()
	}

	printReport()

}
