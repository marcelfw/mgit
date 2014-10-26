// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package command implements all internal commands.
// This source returns the last matching root directory of a repository search string.
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

func (cmd cmdPath) Init(args []string) (outCmd interface{}) {
	if len(args) >= 1 {
		cmd.match = args[0]
	}
	return cmd
}

func (cmd cmdPath) RunConcurrently() bool {
	return true
}

func (cmd cmdPath) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	if repository.PathMatch(cmd.match) {
		return repository, true
	}
	return repository, false
}

func (cmd cmdPath) OutputHeader() []string {
	return nil
}

func (cmd cmdPath) Output(repository repository.Repository) interface{} {
	columns := make([]string, 1, 1)

	columns[0] = repository.GetPath()

	return columns
}
