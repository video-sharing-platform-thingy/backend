package src

import "log"

func check(err error) {
	if err != nil {
		log.Fatalf("Error occured: %s", err)
	}
}
