// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns help.
package command

import (
	"github.com/marcelfw/mgit/repository"
)

type cmdHelp struct {
	command string
}

func NewHelpCommand() cmdHelp {
	var cmd cmdHelp

	return cmd
}

func (cmd cmdHelp) Usage() string {
	return "Show this help information."
}

func (cmd cmdHelp) Help() string {
	if cmd.command != "" {
		return "Showing help about " + cmd.command
	}
	return `Show help information.

Add command as argument to help for more information about the command.`
}

func (cmd cmdHelp) Init(args []string) (outCmd repository.Command) {
	if len(args) >= 1 {
		cmd.command = args[0]
		return cmd
	}
	return nil
}

func (cmd cmdHelp) Output(commands map[string]repository.Command, version string) string {
	if helpCommand, ok := commands[cmd.command]; ok == true {
		return helpCommand.Help()
	}
	return "Unknown command \"" + cmd.command + "\"."
}
