// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repository implements detection, filtering and structure of repositories.
// This source detects and filters repositories.
package repository

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// RepositoryFilter defines a filter for repositories.
type RepositoryFilter struct {
	rootDirectory string
	depth         int

	branch   string
	nobranch bool
	remote   string
	noremote bool
}

// create a new RepositoryFilter
func NewRepositoryFilter(config map[string]string) (filter RepositoryFilter) {
	filter.rootDirectory = "."
	if value, ok := config["rootDirectory"]; ok {
		filter.rootDirectory = value
	}
	if value, ok := config["depth"]; ok {
		depth, err := strconv.ParseInt(value, 10, 0)
		if err == nil {
			filter.depth = int(depth)
		} else {
			filter.depth = 0
		}
	}

	if value, ok := config["branch"]; ok {
		filter.branch = value
	}
	if value, ok := config["nobranch"]; ok {
		filter.branch = value
		filter.nobranch = true
	}
	if value, ok := config["remote"]; ok {
		filter.remote = value
	}
	if value, ok := config["noremote"]; ok {
		filter.remote = value
		filter.noremote = true
	}

	return filter
}

// analysePath extracts repositories from regular file paths.
func analysePath(filter RepositoryFilter, reposChannel chan Repository) filepath.WalkFunc {
	no_of_repositories := 0

	return func(vpath string, f os.FileInfo, err error) error {
		base := path.Base(vpath)
		if base == ".git" {
			// Name is Git-directory without rootDirectory.
			name := strings.TrimLeft(path.Dir(vpath)[len(filter.rootDirectory):], "/")
			if filter.depth > 0 && (strings.Count(name, "/")+1) > filter.depth {
				// if depth limit is set, ignore directories too deep.
				return nil
			}
			if repository, ok := NewRepository(no_of_repositories, name, vpath); ok {
				var found = true
				if found == true && filter.branch != "" {
					is_branch := repository.IsBranch(filter.branch)
					found = (!filter.nobranch && is_branch) || (filter.nobranch && !is_branch)
				}
				if found == true && filter.remote != "" {
					is_remote := repository.IsRemote(filter.remote)
					found = (!filter.noremote && is_remote) || (filter.noremote && !is_remote)
				}

				if found {
					no_of_repositories++
					reposChannel <- repository
				}
			}
		}
		return nil
	}
}

// findRepositories finds and filters repositories below the rootDirectory.
func FindRepositories(filter RepositoryFilter, numDigesters int) chan Repository {
	reposChannel := make(chan Repository, numDigesters)

	go func() {
		filepath.Walk(filter.rootDirectory, analysePath(filter, reposChannel))

		close(reposChannel)
	}()

	return reposChannel
}
