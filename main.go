package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

func main() {
	cli := cli.NewCLI("mirror", Version)
	cli.Args = os.Args[1:]
	cli.Commands = Commands

	exitStatus, err := cli.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)
}
