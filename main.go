// hostess is command-line utility for managing your /etc/hosts file. Works on
// Unixes and Windows.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cbednarski/hostess/hostess"
)

const help = `An idempotent tool for managing %s

Commands

    fmt                  Reformat the hosts file

    add <hostname> <ip>  Add or overwrite a hosts entry
    rm <hostname>        Remote a hosts entry
    on <hostname>        Enable a hosts entry
    off <hostname>       Disable a hosts entry

    ls                   List hosts entries
    has                  Exit 0 if entry present in hosts file, 1 if not

    dump                 Export hosts entries as JSON
    apply                Import hosts entries from JSON

    All commands that change the hosts file will implicitly reformat it.

Flags

    -n will preview changes but not rewrite your hosts file

Configuration

    HOSTESS_FMT may be set to unix or windows to force that platform's syntax
    HOSTESS_PATH may be set to point to a file other than the platform default

About

    Copyright 2015-2020 Chris Bednarski <chris@cbednarski.com>; MIT Licensed
    Portions Copyright the Go authors, licensed under BSD-style license
    Bugs and updates via https://github.com/cbednarski/hostess
`

var (
	Version           = "dev"
	ErrInvalidCommand = errors.New("invalid command")
)

func ExitWithError(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}

func Usage() {
	fmt.Print(help, hostess.GetHostsPath())
	os.Exit(0)
}

func CommandUsage(command string) error {
	return fmt.Errorf("Usage: %s %s <hostname>", os.Args[0], command)
}

func wrappedMain(args []string) error {
	cli := flag.NewFlagSet(args[0], flag.ExitOnError)
	preview := cli.Bool("n", false, "preview")
	cli.Usage = Usage

	command := ""
	if len(args) > 1 {
		command = args[1]
	} else {
		Usage()
	}

	if err := cli.Parse(args[2:]); err != nil {
		return err
	}

	options := &Options{
		Preview: *preview,
	}

	switch command {

	case "-v", "--version", "version":
		fmt.Println(Version)
		return nil

	case "", "-h", "--help", "help":
		cli.Usage()
		return nil

	case "fmt":
		return Format(options)

	case "add":
		if len(cli.Args()) != 2 {
			return fmt.Errorf("Usage: %s add <hostname> <ip>", cli.Name())
		}
		return Add(options, cli.Arg(0), cli.Arg(1))

	case "rm":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return Remove(options, cli.Arg(0))

	case "on":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return Enable(options, cli.Arg(0))

	case "off":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return Disable(options, cli.Arg(0))

	case "ls":
		return List(options)

	case "has":
		if cli.Arg(0) == "" {
			return CommandUsage(command)
		}
		return Has(options, cli.Arg(0))

	case "dump":
		return Dump(options)

	case "apply":
		if cli.Arg(0) == "" {
			return fmt.Errorf("Usage: %s apply <filename>", args[0])
		}
		return Apply(options, cli.Arg(0))

	default:
		return ErrInvalidCommand
	}
}

func main() {
	ExitWithError(wrappedMain(os.Args))
}
