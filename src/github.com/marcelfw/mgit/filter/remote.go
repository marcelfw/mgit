// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filter implements all internal filters.
// This code filters on the presence of a remote.
package filter

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
	"fmt"
)

type filterRemote struct {
	name string

	remote   *string
	noremote *string
}

func NewRemoteFilter() filterRemote {
	filter := filterRemote{name:"remote"}

	return filter
}

func (filter filterRemote) Usage() string {
	return "Filter on presence of remote."
}

func (filter filterRemote) AddFlags(flags *flag.FlagSet) (repository.Filter) {
	filter.remote = flags.String("remote", "", "select only with this remote")
	filter.noremote = flags.String("noremote", "", "select only without this remote")

	return filter
}

func (filter filterRemote) Dump() string {
	return fmt.Sprintf("remote: remote=%s, noremote=%s", *filter.remote, *filter.noremote)
}

func (filter filterRemote) FilterRepository(repos repository.Repository) bool {
	if *filter.remote != "" {
		return repos.IsRemote(*filter.remote)
	} else if *filter.noremote != "" {
		return !repos.IsRemote(*filter.noremote)
	}

	return true
}