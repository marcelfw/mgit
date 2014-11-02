// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns the matching directory of a repository search string.
package command

import (
	"github.com/marcelfw/mgit/repository"
)

type cmdPath struct {
	match string

	repository repository.Repository
}

func NewPathCommand() cmdPath {
	var cmd cmdPath

	return cmd
}

func (cmd cmdPath) Usage() string {
	return "Return repository path of matched names."
}

func (cmd cmdPath) Help() string {
	return `Return repository path of matched names.

Returns all the actual paths matched by the search argument.`
}

func (cmd cmdPath) Init(args []string) (outCmd repository.Command) {
	if len(args) >= 1 {
		cmd.match = args[0]
	}
	return cmd
}

func (cmd cmdPath) IsInteractive() bool {
	return false
}

func (cmd cmdPath) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	if repository.PathMatch(cmd.match) {
		return repository, true
	}
	return repository, false
}

func (cmd cmdPath) Header() string {
	return ""
}

func (cmd cmdPath) Footer() string {
	return ""
}

func (cmd cmdPath) Output(repository repository.Repository) string {
	return repository.GetPath()
}
