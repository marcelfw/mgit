// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package commands implements all internal commands.
// This source returns the status of all repositories.
package commands

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
)

const nameStatus string = "status"

type cmdStatus struct {
}

func NewStatusCommand() cmdStatus {
	var cmd cmdStatus

	return cmd
}

func (cmd cmdStatus) Usage(name_len int) string {
	return nameStatus + strings.Repeat(" ", name_len-len(nameStatus)) + " Return the status of each repository"
}

func (cmd cmdStatus) Help() string {
	return "Returns really short status for repository."
}

func (cmd cmdStatus) Init(args []string) (outCmd interface{}) {
	// we don't do anything
	return nil

}

func (cmd cmdStatus) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	// we require what we already have
	return repository, true
}

func (cmd cmdStatus) OutputHeader() []string {
	return []string{
		"Name", "Branch", "Status",
	}
}

func (cmd cmdStatus) Output(repository repository.Repository) []string {
	columns := make([]string, 3, 3)

	columns[0] = repository.Name
	columns[1] = repository.GetCurrentBranch()
	columns[2] = repository.GetStatusJudgement()

	return columns
}
