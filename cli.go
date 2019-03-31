package main

import (
	"archive/zip"
	"io"
	"log"
	"os"

	"io/ioutil"
	"net/http"
	"path/filepath"
)

var starterRepo = "https://github.com/pulsar-go/example/archive/master.zip"

func main() {
	// argument 'ǹame'
	name := os.Args[1]
	// (option 'version')

	path, err := filepath.Abs(name)
	exitOnError(err)

	// if cwd/name does exist => exit
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Fatal(err)
	}

	// create directory in cwd() . /name
	err = os.MkdirAll(path, 0755)
	exitOnError(err)

	// (get version)

	// download project-zip
	// https://github.com/pulsar-go/example/archive/master.zip
	response, err := http.Get(starterRepo)
	exitOnError(err)
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	exitOnError(err)

	file, err := os.Create(path + "/tmp")
	exitOnError(err)
	defer file.Close()

	_, err = file.Write(data)
	exitOnError(err)

	// extract zip
	var filenames []string

	reader, err := zip.OpenReader(path + "/tmp")
	exitOnError(err)
	defer reader.Close()

	for _, file := range reader.File {
		fpath := filepath.Join(path, file.Name)

		filenames = append(filenames, fpath)

		if file.FileInfo().IsDir() {
			// Make folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		exitOnError(err)

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		exitOnError(err)

		rc, err := file.Open()
		exitOnError(err)

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		exitOnError(err)
	}

	// files, err := ioutil.ReadDir("./example-master")
	// exitOnError(err)

	// for _, file := range files {
	// 	oldLocation := filepath.Join(path+"/example-master", file.Name)
	// 	newLocation := filepath.Join(path, file.Name)
	// 	err := os.Rename(oldLocation, newLocation)
	// 	exitOnError(err)
	// }

	// delete zip
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
