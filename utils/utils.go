package utils

import "log"

func Check(err error) {
	if err != nil {
		log.Fatalf("Error occured: %s", err)
	}
}
