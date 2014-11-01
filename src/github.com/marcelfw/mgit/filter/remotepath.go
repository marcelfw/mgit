// Copyright (c) 2014 Marcel Wouters

// Package filter implements all internal filters.
// This code filters on the presence of a remote.
package filter

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
	"fmt"
)

type filterRemotePath struct {
	name string

	remotepath   *string
	noremotepath *string
}

func NewRemotePathFilter() filterRemotePath {
	filter := filterRemotePath{name:"remotepath"}

	return filter
}

func (filter filterRemotePath) Usage() string {
	return "Filter on a remote path."
}

func (filter filterRemotePath) AddFlags(flags *flag.FlagSet) (repository.Filter) {
	filter.remotepath = flags.String("remotepath", "", "select only with this value is found in the remote path")
	filter.noremotepath = flags.String("noremotepath", "", "select only when this value is not found in the remote path")

	return filter
}

func (filter filterRemotePath) Dump() string {
	return fmt.Sprintf("remotepath: remotepath=%s, noremotepath=%s", *filter.remotepath, *filter.noremotepath)
}

func (filter filterRemotePath) FilterRepository(repos repository.Repository) bool {
	if *filter.remotepath != "" {
		if !repos.RemotePathContains(*filter.remotepath) {
			return false
		}
	}
	if *filter.noremotepath != "" {
		if repos.RemotePathContains(*filter.noremotepath) {
			return false
		}
	}

	return true
}
