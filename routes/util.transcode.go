package routes

import (
	"log"
	"time"

	"github.com/xfrr/goffmpeg/transcoder"
)

func transcode(fromPath string, toPath string) error {
	start := time.Now()
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(fromPath, toPath)
	if err != nil {
		return err
	}
	done := trans.Run(false)
	err = <-done
	if err != nil {
		return err
	}
	t := time.Now()
	elapsed := t.Sub(start)
	log.Println("Transcoded in", elapsed)
	return nil
}
