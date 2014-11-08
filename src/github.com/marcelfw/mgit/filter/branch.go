// Copyright (c) 2014 Marcel Wouters

// Package filter implements all internal filters.
// This code filters on the presence of a branch.
package filter

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
	"io/ioutil"
	"log"
	"os"
)

type filterBranch struct {
	name string

	branch   *string
	nobranch *string
}

func NewBranchFilter() filterBranch {
	filter := filterBranch{name: "branch"}

	return filter
}

func (filter filterBranch) Usage() string {
	return "Filter on the present of a branch."
}

func (filter filterBranch) AddFlags(flags *flag.FlagSet) repository.Filter {
	filter.branch = flags.String("branch", "", "select only with this branch")
	filter.nobranch = flags.String("nobranch", "", "select only without this branch")

	return filter
}

// getBranches returns the branches.
func getBranches(repository repository.Repository) (branches map[string]bool) {
	branches = make(map[string]bool)

	if fi, err := os.Stat(repository.GetGitRoot() + "/refs/heads"); err == nil && fi.IsDir() {
		if fis, err := ioutil.ReadDir(repository.GetGitRoot() + "/refs/heads"); err == nil {
			for _, fi := range fis {
				// We don't support branches in subdirectories.
				if !fi.IsDir() {
					branches[fi.Name()] = true
				}
			}
		}
	} else {
		log.Printf("! no directory [%v]", err)
	}

	//log.Printf("Branches for repository \"%s\" => \"%v\"", repository.Name, branches)

	return branches
}

func (filter filterBranch) FilterRepository(repos repository.Repository) bool {
	branches := getBranches(repos)

	if *filter.branch != "" {
		if _, ok := branches[*filter.branch]; !ok {
			return false
		}
	}
	if *filter.nobranch != "" {
		if _, ok := branches[*filter.nobranch]; ok {
			return false
		}
	}

	return true
}
