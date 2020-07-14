package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Counter
var count uint = 0
var processedCount uint = 0

// STATUS Type for file progress
type STATUS uint

// Status for file process progress
const (
	YETTOSTART STATUS = 0
	INPROGRESS STATUS = 1
	DONE       STATUS = 3
	FAILED     STATUS = 4
)

// MODE is a mode used to describe parallel or serial
type MODE uint

// execution Modes
const (
	PARALLEL MODE = 0
	SERIAL   MODE = 1
)

const (
	//DEFAULTNOFILES is the default no of files to be processed simultaneously
	DEFAULTNOFILES uint = 20
)

// StatusStringMap is used to create a lookup for status
var StatusStringMap = map[STATUS]string{YETTOSTART: "NOTSTARTED", INPROGRESS: "INPROGRESS", DONE: "DONE", FAILED: "FAILED"}

// AviFiles is a data structure used for storing information
type AviFiles struct {
	inFilepath      string
	outFilename     string
	outputDir       string
	no              uint
	status          STATUS
	mtx             sync.Mutex
	processSignal   chan bool
	printStatusOnce bool
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

		count++

		aviobj = &AviFiles{
			inFilepath:      path,
			outFilename:     getOutFilename(path),
			outputDir:       getOutputDir(path),
			no:              count,
			status:          YETTOSTART,
			processSignal:   make(chan bool),
			printStatusOnce: false,
		}

		log.Printf("aviobj: inFilepath: %s, outFilename: %s, outputDir: %s, no: %d\n\n",
			aviobj.inFilepath, aviobj.outFilename, aviobj.outputDir, aviobj.no)
	}

	if aviobj != nil {
		ProcessingQueue = append(ProcessingQueue, aviobj)
	}

	return nil
}

func main() {
	// Set log file
	logFile, err := os.OpenFile("gomediacomverter.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("\nFailed to create logfile, error: %s\n", err.Error())
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	// Set input parameters or flags
	inputDir := flag.String("inputdir", "", "The input directory where avi files are stored. All files under this folder will be Recursively processed")
	dryRun := flag.Bool("dryrun", false, "Only list the files to be processed")
	noFiles := flag.Uint("nofiles", DEFAULTNOFILES, "no of files to process simulatenosuly")
	serialMode := flag.Bool("serial", true, "All files to be processed serially")
	parallelMode := flag.Bool("parallel", false, "All files to be processed simulatenosuly")

	flag.Parse()

	// validate flags
	if *inputDir == "" {
		log.Printf("Missing input Directory arg")
		exit("Missing input Directory arg", -1)
	}

	if *noFiles == 0 {
		exit("You choose no files to process..increase default or overide with a meaningful nooffiles", 0)
	}

	if _, err := os.Stat(*inputDir); os.IsNotExist(err) {
		exit("directory does not exit", 2)
	}

	if *serialMode || *parallelMode {

	} else {
		exit("Cant run in serial mode and parallel at the same time", 3)
	}

	if *serialMode {
		fmt.Println("Running in serial mode")
	} else {
		fmt.Println("Running in parallel mode")
	}

	filepath.Walk(*inputDir, mediawalk)

	if count == 0 {
		exit("There are no files to be processed", 0)
	} else {
		fmt.Printf("\nThere are in total %d eligible files to be processed in this folder\n", count)

		var noFilesMessage string = ""
		if count < *noFiles {
			*noFiles = count
			noFilesMessage = fmt.Sprintf("%d files", *noFiles)
		} else if *noFiles == DEFAULTNOFILES {
			noFilesMessage = fmt.Sprintf("(Default of %d files)", DEFAULTNOFILES)
		} else {
			noFilesMessage = fmt.Sprintf("%d files", *noFiles)
		}

		fmt.Printf("\n%s will be processed\n", noFilesMessage)
	}

	fmt.Println("*************************************")
	fmt.Println(" Eligible Avi Files to be Processed")
	fmt.Println("************************************")
	for j := range ProcessingQueue {
		fmt.Printf("%s\n", ProcessingQueue[j].inFilepath)
	}

	if *dryRun {
		exit("", 0)
	}

	var wg sync.WaitGroup
	fmt.Println()
	fmt.Println("********")
	fmt.Println("Progress")
	fmt.Println("********")

	for i := uint(0); i < *noFiles; i++ {
		if *parallelMode {
			wg.Add(1)
			go ProcessingQueue[i].runConversion(PARALLEL, &wg)
		} else {
			ProcessingQueue[i].runConversion(SERIAL, &wg)
			fmt.Printf("(%d/%d). %s\t ...... [%s]\n", ProcessingQueue[i].no, count, ProcessingQueue[i].inFilepath, StatusStringMap[ProcessingQueue[i].status])
		}
	}

	if *parallelMode {
		var countmtx sync.Mutex
		for {
			for i := uint(0); i < *noFiles; i++ {
				select {
				case <-ProcessingQueue[i].processSignal:
					countmtx.Lock()
					if !ProcessingQueue[i].printStatusOnce {
						fmt.Printf("(%d/%d). %s\t ...... [%s]\n", ProcessingQueue[i].no, count, ProcessingQueue[i].inFilepath, StatusStringMap[ProcessingQueue[i].status])
						ProcessingQueue[i].printStatusOnce = true
						processedCount++
					}
					countmtx.Unlock()
				case <-time.After(30 * time.Second):
				}
			}

			if processedCount == *noFiles {
				break
			}
		}
		wg.Wait()
	}

}
