// Copyright (c) 2014 Marcel Wouters

// Package filter implements all internal filters.
// This code filters on the name of the repository.
package filter

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
)

type filterName struct {
	name string

	match *string
}

func NewNameFilter() filterName {
	filter := filterName{name: "name"}

	return filter
}

func (filter filterName) Name() string {
	return filter.name
}

func (filter filterName) Usage() map[string]string {
	return map[string]string{
		"-name <partial-name>": "Match on partial name match.",
	}
}

func (filter filterName) AddFlags(flags *flag.FlagSet) repository.Filter {
	filter.match = flags.String("name", "", "select only when name is found")

	return filter
}

func (filter filterName) FilterRepository(repos repository.Repository) bool {
	if *filter.match != "" {
		if !repos.NameContains(*filter.match) {
			return false
		}
	}

	return true
}
