package main

import (
	"github.com/mefellows/mirror/command"
	"github.com/mitchellh/cli"
	"os"
)

var Commands map[string]cli.CommandFactory
var Ui cli.Ui

func init() {

	Ui = &cli.ColoredUi{
		Ui:          &cli.BasicUi{Writer: os.Stdout, Reader: os.Stdin, ErrorWriter: os.Stderr},
		OutputColor: cli.UiColorYellow,
		InfoColor:   cli.UiColorNone,
		ErrorColor:  cli.UiColorRed,
	}

	meta := command.Meta{
		Ui: Ui,
	}

	Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &command.SyncCommand{
				Meta: meta,
			}, nil
		},
	}
}
