package main

import (
	"os"

	"github.com/pulsar-go/cli/commands"
	"github.com/pulsar-go/cli/utils"
)

var commandList = []*commands.Command{
	commands.NewCommand,
}

func main() {
	commandName, args := parseArgs()

	var foundCommand *commands.Command

	for i := range commandList {
		if commandList[i].Name == commandName {
			foundCommand = commandList[i]
		}
	}

	if foundCommand == nil {
		utils.ExitOnNewError("I don't know that command: " + commandName)
	}

	foundCommand.Action(args...)
}

func parseArgs() (commandName string, args []string) {
	return os.Args[1], os.Args[2:]
}
