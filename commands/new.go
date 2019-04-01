package commands

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pulsar-go/cli/utils"
)

const (
	// PulsarFramework Pulsar framework package reference
	PulsarFramework = "github.com/pulsar-go/pulsar"

	// SkeletonRepo Pulsar skeleton repo to start new projects
	SkeletonRepo = "https://github.com/pulsar-go/example/archive/master.zip"
)

// NewCommand creates a new pulsar project
var NewCommand = &Command{
	Name:        "new",
	Description: "creates a new pulsar project",
	Action: func(args ...string) {
		run(args[0])
	},
}

type app struct {
	name  string
	data  []byte
	paths map[string]string
}

func newApp() *app {
	a := &app{}
	a.paths = map[string]string{
		"base":     "",
		"zip":      "",
		"unzipped": "",
	}

	return a
}

func run(name string) {
	cmd := newApp()

	cmd.setName(name).
		setPath().
		verifyItDoesntExist().
		createDir().
		download().
		extract().
		moveUp().
		clean().
		install()
}

// set the app's name
func (cmd *app) setName(name string) *app {
	cmd.name = name

	return cmd
}

// set the app's path
func (cmd *app) setPath() *app {
	path, err := filepath.Abs(cmd.name)
	utils.ExitOnError(err)

	cmd.paths["base"] = path
	return cmd
}

// verify app was not already created
func (cmd *app) verifyItDoesntExist() *app {
	// if cwd/name does exist => exit
	if _, err := os.Stat(cmd.paths["base"]); !os.IsNotExist(err) {
		utils.ExitOnNewError("Folder with name " + cmd.name + " already exists!")
	}

	return cmd
}

// create app directory
func (cmd *app) createDir() *app {
	// create directory in cwd() . /name
	err := os.MkdirAll(cmd.paths["base"], os.ModePerm)
	utils.ExitOnError(err)

	return cmd
}

// download app skeleton and save the zip file
func (cmd *app) download() *app {
	// download project-zipPulsar skeleton repo to start new projects
	response, err := http.Get(SkeletonRepo)
	utils.ExitOnError(err)
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	utils.ExitOnError(err)

	cmd.data = data

	return cmd.writeZip()
}

// write zip file to file system
func (cmd *app) writeZip() *app {
	cmd.paths["zip"] = cmd.paths["base"] + "/tmp"
	file, err := os.Create(cmd.paths["zip"])
	utils.ExitOnError(err)
	defer file.Close()

	_, err = file.Write(cmd.data)
	utils.ExitOnError(err)

	return cmd
}

// extract zip file
func (cmd *app) extract() *app {
	// extract zip
	reader, err := zip.OpenReader(cmd.paths["zip"])
	utils.ExitOnError(err)
	defer reader.Close()

	for _, file := range reader.File {
		fpath := filepath.Join(cmd.paths["base"], file.Name)

		if file.FileInfo().IsDir() {
			// Make folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
		utils.ExitOnError(err)

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		utils.ExitOnError(err)

		rc, err := file.Open()
		utils.ExitOnError(err)

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		utils.ExitOnError(err)
	}

	return cmd
}

// move all files & folders one layer up
func (cmd *app) moveUp() *app {
	cmd.paths["unzipped"] = cmd.paths["base"] + "/example-master"
	files, err := ioutil.ReadDir(cmd.paths["unzipped"])
	utils.ExitOnError(err)

	for _, file := range files {
		oldLocation := filepath.Join(cmd.paths["unzipped"], file.Name())
		newLocation := filepath.Join(cmd.paths["base"], file.Name())
		err := os.Rename(oldLocation, newLocation)
		utils.ExitOnError(err)
	}

	return cmd
}

// delete downloaded zip and empty unzipped directory
func (cmd *app) clean() *app {
	elements := [2]string{cmd.paths["unzipped"], cmd.paths["zip"]}

	for _, element := range elements {
		os.Remove(element)
	}

	return cmd
}

// install pulsar framework if not already installed
func (cmd *app) install() {
	// install PulsarFramework if not already installed
	_, err := exec.Command("sh", "-c", "go list "+PulsarFramework).Output()

	if err != nil {
		_, err := exec.Command("sh", "-c", "go get "+PulsarFramework).Output()
		utils.ExitOnError(err)
	}
}
