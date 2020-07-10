package main

import (
	"fmt"
	"os/exec"
)

func (queue *AviFiles) runConversion() error {

	queue.mtx.Lock()

	args := []string{
		"-i",
		queue.inFilepath,
		queue.outFilename,
	}

	if queue.getStatus() != YETTOSTART {
		fmt.Printf("status is %d", queue.getStatus())
		return nil
	}

	queue.setStatus(INPROGRESS)
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.Output()

	if err != nil {
		fmt.Printf("Problem Converting file:%s, error: %s", queue.getInputFile(), err.Error())
		queue.setStatus(FAILED)
	} else {
		queue.setStatus(DONE)
	}

	queue.mtx.Unlock()

	return nil

}
