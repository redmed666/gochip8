package util

import (
	"log"
)

// CheckError check if there is an error
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
