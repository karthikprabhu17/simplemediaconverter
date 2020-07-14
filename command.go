package main

import (
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func (queue *AviFiles) runConversion(mode MODE, wg *sync.WaitGroup) error {
	if mode == PARALLEL {
		defer wg.Done()
	}

	queue.mtx.Lock()
	defer queue.mtx.Unlock()

	args := []string{
		"-i",
		queue.inFilepath,
		queue.outFilename,
	}

	if queue.getStatus() != YETTOSTART {
		fmt.Printf("status is %d", queue.getStatus())
		if mode == PARALLEL {
			queue.processSignal <- true
			close(queue.processSignal)
		}
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

	if mode == PARALLEL {
		queue.processSignal <- true
		close(queue.processSignal)
	}

	return nil

}
