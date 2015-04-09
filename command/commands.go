package command

import (
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

	meta := Meta{
		Ui: Ui,
	}

	Commands = map[string]cli.CommandFactory{
		"daemon": func() (cli.Command, error) {
			return &DaemonCommand{
				Meta: meta,
			}, nil
		},
		"remote": func() (cli.Command, error) {
			return &RemoteCommand{
				Meta: meta,
			}, nil
		},
		"sync": func() (cli.Command, error) {
			return &SyncCommand{
				Meta: meta,
			}, nil
		},
		"pki": func() (cli.Command, error) {
			return &PkiCommand{
				Meta: meta,
			}, nil
		},
	}
}
