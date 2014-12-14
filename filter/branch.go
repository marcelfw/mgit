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
	"strings"
)

type filterBranch struct {
	name string

	current  *string
	branch   *string
	nobranch *string
}

func NewBranchFilter() filterBranch {
	filter := filterBranch{name: "branch"}

	return filter
}

func (filter filterBranch) Name() string {
	return filter.name
}

func (filter filterBranch) Usage() map[string]string {
	return map[string]string{
		"-branch <branch>":   "Match when branch is found.",
		"-nobranch <branch>": "Match only when branch is not found.",
	}
}

func (filter filterBranch) AddFlags(flags *flag.FlagSet) repository.Filter {
	filter.current = flags.String("current", "", "select only when current matches")
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

	if *filter.current != "" {
		if branch, _, ok := repos.ExecGit("rev-parse", "--abbrev-ref", "HEAD"); ok {
			if *filter.current != strings.TrimRight(branch, "\r\n") {
				return false
			}
		}
	}

	if *filter.branch != "" {
		if *filter.branch == "master" {
			// there might not be any refs yet for "master"
			// master is always assumed to be there
			return true
		}
		if _, ok := branches[*filter.branch]; !ok {
			return false
		}
	}
	if *filter.nobranch != "" {
		if *filter.nobranch == "master" {
			// master is always assumed to be there
			// (so this basically filters everything away)
			return false
		}
		if _, ok := branches[*filter.nobranch]; ok {
			return false
		}
	}

	return true
}
