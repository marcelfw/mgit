// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package command implements all internal commands.
// This source returns the last matching root directory of a repository search string.
package command

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
)

type cmdEcho struct {
	args []string

	repository repository.Repository
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

Pares the rest of the command-line and performs all macro conversion.
Useful for testing macros.`
}

func (cmd cmdEcho) Init(args []string) (outCmd repository.Command) {
	cmd.args = args
	return cmd
}

func (cmd cmdEcho) IsInteractive() bool {
	return false
}

func (cmd cmdEcho) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	return repository, true
}

func (cmd cmdEcho) OutputHeader() []string {
	return nil
}

func (cmd cmdEcho) Output(repository repository.Repository) interface{} {
	return strings.Join(repository.ReplaceMacros(cmd.args), " ")
}
