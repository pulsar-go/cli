package main

import (
	"errors"
	"io"
	"log"
	"os"

	"archive/zip"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
)

const (
	// PulsarFramework Pulsar framework package reference
	PulsarFramework = "github.com/pulsar-go/pulsar"

	// SkeletonRepo Pulsar skeleton repo to start new projects
	SkeletonRepo = "https://github.com/pulsar-go/example/archive/master.zip"
)

func main() {
	// argument 'ǹame'
	name := os.Args[1]

	path, err := filepath.Abs(name)
	exitOnError(err)

	// if cwd/name does exist => exit
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Fatal(errors.New("Folder with name " + name + " already exists!"))
	}

	// create directory in cwd() . /name
	err = os.MkdirAll(path, 0755)
	exitOnError(err)

	// download project-zipPulsar skeleton repo to start new projects
	response, err := http.Get(SkeletonRepo)
	exitOnError(err)
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	exitOnError(err)

	zipPath := path + "/tmp"
	file, err := os.Create(zipPath)
	exitOnError(err)
	defer file.Close()

	_, err = file.Write(data)
	exitOnError(err)

	// extract zip
	reader, err := zip.OpenReader(zipPath)
	exitOnError(err)
	defer reader.Close()

	for _, file := range reader.File {
		fpath := filepath.Join(path, file.Name)

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

	unzippedFolder := path + "/example-master"
	files, err := ioutil.ReadDir(unzippedFolder)
	exitOnError(err)

	for _, file := range files {
		oldLocation := filepath.Join(unzippedFolder, file.Name())
		newLocation := filepath.Join(path, file.Name())
		err := os.Rename(oldLocation, newLocation)
		exitOnError(err)
	}

	// delete zip
	os.Remove(unzippedFolder)
	os.Remove(zipPath)

	// install PulsarFramework if not already installed
	_, err = exec.Command("sh", "-c", "go list "+PulsarFramework).Output()
	if err != nil {
		_, err := exec.Command("sh", "-c", "go get "+PulsarFramework).Output()
		exitOnError(err)
	}
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
