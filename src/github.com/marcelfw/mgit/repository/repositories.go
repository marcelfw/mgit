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
			name := strings.TrimPrefix(path.Dir(vpath), filter.rootDirectory)
			name = strings.TrimLeft(name, "/")
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
