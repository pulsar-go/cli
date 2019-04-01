package utils

import (
	"errors"
	"log"
)

// ExitOnError log err and exits app
func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ExitOnNewError creates a new error from string and logs it
func ExitOnNewError(message string) {
	ExitOnError(errors.New(message))
}
