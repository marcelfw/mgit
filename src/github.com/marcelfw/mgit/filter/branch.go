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
