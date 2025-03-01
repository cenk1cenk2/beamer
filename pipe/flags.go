package pipe

import (
	"github.com/urfave/cli/v2"
	. "gitlab.kilic.dev/libraries/plumber/v5"
)

//revive:disable:line-length-limit

const (
	CATEGORY_CONFIG = "Config"
)

var Flags = []cli.Flag{
	&cli.StringFlag{
		Category:    CATEGORY_CONFIG,
		Name:        "config-file",
		Usage:       "Configuration file to read from. json(https://raw.githubusercontent.com/cenk1cenk2/docker-vizier/main/schema.json)",
		Required:    false,
		Value:       "",
		EnvVars:     []string{"VIZIER_CONFIG_FILE"},
		Destination: &TL.Pipe.File,
	},
}

//revive:disable:unused-parameter
func ProcessFlags(tl *TaskList[Pipe]) error {
	return nil
}
