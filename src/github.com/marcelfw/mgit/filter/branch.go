// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package commands implements all internal commands.
// This source returns the status of all repositories.
package filter

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
)

type filterBranch struct {
	name string

	branch   string
	nobranch string
}

func NewBranchFilter() filterBranch {
	filter := filterBranch{name:"branch"}

	return filter
}

func (cmd filterBranch) Usage() string {
	return "Filter on branch."
}

func (cmd filterBranch) AddFlags(flags *flag.FlagSet) {
	flags.StringVar(&cmd.branch, "b", "", "select only with this branch")
	flags.StringVar(&cmd.nobranch, "nb", "", "select only without this branch")
}

func (cmd filterBranch) FilterRepository(repository.Repository) bool {
	return true
}
