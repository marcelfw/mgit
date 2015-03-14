// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns the version of the code.
package command

import (
	"github.com/marcelfw/mgit/repository"
)

type cmdVersion struct {
}

func NewVersionCommand() cmdVersion {
	var cmd cmdVersion

	return cmd
}

func (cmd cmdVersion) Usage() string {
	return "Show current version."
}

func (cmd cmdVersion) Help() string {
	return "Show current version."
}

func (cmd cmdVersion) Init(args []string, interactive bool, dryrun bool) (outCmd repository.Command) {
	return nil
}

func (cmd cmdVersion) Output(commands map[string]repository.Command, version string) string {
	return "mgit version " + version
}
