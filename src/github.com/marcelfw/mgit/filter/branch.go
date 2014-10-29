// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filter implements all internal filters.
// This code filters on the presence of a branch.
package filter

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
	"fmt"
)

type filterBranch struct {
	name string

	branch   *string
	nobranch *string
}

func NewBranchFilter() filterBranch {
	filter := filterBranch{name:"branch"}

	return filter
}

func (filter filterBranch) Usage() string {
	return "Filter on the present of a branch."
}

func (filter filterBranch) AddFlags(flags *flag.FlagSet) (repository.Filter) {
	filter.branch = flags.String("branch", "", "select only with this branch")
	filter.nobranch = flags.String("nobranch", "", "select only without this branch")

	return filter
}

func (filter filterBranch) Dump() string {
	return fmt.Sprintf("branch: branch=%s, nobranch=%s", *filter.branch, *filter.nobranch)
}

func (filter filterBranch) FilterRepository(repos repository.Repository) bool {
	if *filter.branch != "" {
		if repos.IsBranch(*filter.branch) {
			return true
		}
	}
	if *filter.nobranch != "" {
		return !repos.IsBranch(*filter.nobranch)
	}
	
	return true
}
