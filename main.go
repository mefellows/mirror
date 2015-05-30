package main

import (
	"fmt"
	"github.com/mefellows/mirror/command"
	_ "github.com/mefellows/mirror/filesystem/fs"
	_ "github.com/mefellows/mirror/filesystem/remote"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	cli := cli.NewCLI("mirror", Version)
	cli.Args = os.Args[1:]
	cli.Commands = command.Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(exitStatus)
}
