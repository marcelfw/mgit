// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package command implements all internal commands.
// This source returns help.
package command

import (
	"github.com/marcelfw/mgit/repository"
)

type cmdHelp struct {
	IsHelp bool
}

func NewHelpCommand() cmdHelp {
	var cmd cmdHelp

	cmd.IsHelp = true

	return cmd
}

func (cmd cmdHelp) Usage() string {
	return "Show this help information."
}

func (cmd cmdHelp) Help() string {
	return `Show help information.

Add command as argument to help for more information about the command.`
}

func (cmd cmdHelp) Init(args []string) (outCmd interface{}) {
	return nil
}

func (cmd cmdHelp) RunConcurrently() bool {
	return true
}

func (cmd cmdHelp) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	return repository, false
}

func (cmd cmdHelp) OutputHeader() []string {
	return nil
}

func (cmd cmdHelp) Output(repository repository.Repository) interface{} {
	return nil
}
