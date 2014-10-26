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
	"strings"
)

// RepositoryFilter defines a filter for repositories.
type RepositoryFilter struct {
	rootDirectory string
	depth         int

	filters []Filter
}

// create a new RepositoryFilter
func NewRepositoryFilter(rootDirectory string, depth int, filters []Filter) (filter RepositoryFilter) {
	filter.rootDirectory = rootDirectory
	filter.depth = depth

	filter.filters = filters

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
				for _, filter := range filter.filters {
					if filter.FilterRepository(repository) == false {
						found = false
						break
					}
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
