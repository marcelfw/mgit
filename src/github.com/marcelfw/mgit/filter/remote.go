// Copyright (c) 2014 Marcel Wouters
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS
// OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
// OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
		if repos.IsRemote(*filter.remote) {
			return true
		}
	}
	if *filter.noremote != "" {
		return !repos.IsRemote(*filter.noremote)
	}

	return true
}
