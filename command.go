package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func (queue *AviFiles) runConversion(wg *sync.WaitGroup) error {
	defer wg.Done()

	//queue.mtx.Lock()

	args := []string{
		"-i",
		queue.inFilepath,
		queue.outFilename,
	}

	//defer queue.mtx.Unlock()

	if queue.getStatus() != YETTOSTART {
		fmt.Printf("status is %d", queue.getStatus())
		queue.processSignal <- true
		close(queue.processSignal)
		return nil
	}

	queue.setStatus(INPROGRESS)
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.Output()

	if err != nil {
		log.Printf("\nProblem Converting file:%s, error: %s", queue.getInputFile(), err.Error())
		queue.setStatus(FAILED)
	} else {
		queue.setStatus(DONE)
	}

	queue.processSignal <- true
	close(queue.processSignal)
	return nil

}
