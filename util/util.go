package util

import "log"

// CheckError simply will log the error
// and panic if it's not nil.
func CheckError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}
