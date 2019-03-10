package utils

import "log"

// Check exits and logs if there's an error.
func Check(err error) {
	if err != nil {
		log.Fatalf("Error occured: %s", err)
	}
}
