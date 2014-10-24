// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package commands implements all internal commands.
// This source returns the last matching root directory of a repository search string.
package commands

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
)

const namePath string = "path"

type cmdPath struct {
	match string

	repository repository.Repository
}

func NewPathCommand() cmdPath {
	var cmd cmdPath

	return cmd
}

func (cmd cmdPath) Usage(name_len int) string {
	return namePath + strings.Repeat(" ", name_len-len(namePath)) + " Return repository paths of matching names."
}

func (cmd cmdPath) Help() string {
	return "Returns repository directory."
}

func (cmd cmdPath) Init(args []string) (outCmd interface{}) {
	if len(args) >= 1 {
		cmd.match = args[0]
	}
	return cmd
}

func (cmd cmdPath) RunConcurrently() (bool) {
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
	/*return []string{
		"Name", "Branch", "Status",
	}*/
}

func (cmd cmdPath) Output(repository repository.Repository) []string {
	columns := make([]string, 1, 1)

	columns[0] = repository.GetPath()

	return columns
}
