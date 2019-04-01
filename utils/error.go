package utils

import "log"

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
