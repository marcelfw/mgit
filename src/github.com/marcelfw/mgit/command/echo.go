// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns the last matching root directory of a repository search string.
package command

import (
	"strings"

	"github.com/marcelfw/mgit/repository"
)

type cmdEcho struct {
	args []string
}

func NewEchoCommand() cmdEcho {
	var cmd cmdEcho

	return cmd
}

func (cmd cmdEcho) Usage() string {
	return "Echo output after conversion."
}

func (cmd cmdEcho) Help() string {
	return `Echo output after conversion.

Parses the rest of the command-line and performs all macro conversion.
Useful for testing macros.`
}

func (cmd cmdEcho) Init(args []string, interactive bool, dryrun bool) (outCmd repository.Command) {
	cmd.args = args
	return cmd
}

func (cmd cmdEcho) IsInteractive() bool {
	return false
}

func (cmd cmdEcho) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	return repository, true
}

func (cmd cmdEcho) Header() string {
	return ""
}

func (cmd cmdEcho) Footer() string {
	return ""
}

func (cmd cmdEcho) Output(repository repository.Repository) string {
	return strings.Join(repository.ReplaceMacros(cmd.args), " ")
}
