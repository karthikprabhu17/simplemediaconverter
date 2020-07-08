package main

import (
	"fmt"
	"os/exec"
)

func (queue *AviFiles) runConversion() error {

	queue.mtx.Lock()

	queue.setStatus(INPROGRESS)

	args := []string{
		"-i",
		queue.inFilepath,
		queue.outFilename,
	}

	cmd := exec.Command("ffmpeg", args...)
	out, err := cmd.Output()

	if err != nil {
		fmt.Printf("Problem Converting file:%s", queue.getInputFile())
	}

	queue.mtx.Unlock()

	return nil

}
