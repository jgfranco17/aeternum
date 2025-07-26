package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jgfranco17/aeternum/cli/core"
	"github.com/jgfranco17/aeternum/cli/logging"
)

const (
	projectName        = "aeternum"
	projectDescription = "Aeternum: Monitoring REST services with ease."
)

var (
	version string = "0.0.0-dev.1"
)

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&logging.CustomFormatter{})
}

func main() {
	commandsList := []*cobra.Command{
		core.GetPingCommand(),
	}
	command := core.NewCommandRegistry(projectName, projectDescription, version)
	command.RegisterCommands(commandsList)

	err := command.Execute()
	if err != nil {
		log.Error(err.Error())
	}
}
