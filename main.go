package main

import (
	"log"
	"os"

	"net/http"
)

func main() {
	// argument 'ǹame'
	name := os.Args[1]
	// (option 'version')

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := cwd + "/" + name

	// if cwd/name does exist => exit
	if _, err := os.Stat(path); os.IsExist(err) {
		log.Fatal(err)
	}

	// create directory in cwd() . /name
	err = os.MkdirAll(path, 0755)
	if err != nil {
		log.Fatal(err)
	}

	// (get version)

	// download project-zip
	// https://github.com/pulsar-go/example/archive/master.zip
	response, err := http.Get("https://github.com/pulsar-go/example/archive/master.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	// extract zip
	// delete zip
}
