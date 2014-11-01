// Copyright (c) 2014 Marcel Wouters

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
		if !repos.IsBranch(*filter.branch) {
			return false
		}
	}
	if *filter.nobranch != "" {
		if repos.IsBranch(*filter.nobranch) {
			return false
		}
	}
	
	return true
}
