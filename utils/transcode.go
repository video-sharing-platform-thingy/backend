package utils

import (
	"log"
	"time"

	"github.com/xfrr/goffmpeg/transcoder"
)

// Transcode transcodes a video file.
func Transcode(fromPath string, toPath string) error {
	// Get the current time
	start := time.Now()

	// Initialize a new transcoder
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(fromPath, toPath)
	if err != nil {
		return err
	}

	// Run the transcoding
	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}

	// Get the elapsed time and print it
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("Transcoded in", elapsed)

	return nil
}
