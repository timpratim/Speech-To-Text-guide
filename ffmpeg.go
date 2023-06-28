package main

import (
	"fmt"
	"os"
	"os/exec"
)

func ConvertFile(inputFile string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-ac", "1", "-ar", "16000", "-codec:a", "pcm_s16le", "-f", "wav", outputFile)
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return nil
}
