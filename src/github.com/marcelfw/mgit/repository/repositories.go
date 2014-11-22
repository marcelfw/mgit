// Copyright (c) 2014 Marcel Wouters

// Package repository implements detection, filtering and structure of repositories.
// This source detects and filters repositories.
package repository

import (
	"log"
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
		name := ""
		base := path.Base(vpath)
		gitPath := path.Dir(vpath)

		// Name is Git-directory without rootDirectory.
		if gitPath != ".git" {
			name = strings.TrimPrefix(gitPath, filter.rootDirectory)
			name = strings.TrimLeft(name, "/")
			name = strings.TrimSuffix(name, "/.git")
			if filter.depth > 0 && (strings.Count(name, "/")+1) > filter.depth {
				// if depth limit is set, ignore directories too deep.
				log.Printf("Skipping repository \"%s\" (filtered by depth)", name)
				return nil
			}
		}

		var repository Repository
		foundRepository := false
		if base == "HEAD" {
			if _, err := os.Stat(path.Dir(vpath) + "/config"); err == nil {
				repository, foundRepository = NewRepository(no_of_repositories, name, gitPath)
			}
		}

		if foundRepository {
			var allow = true
			for _, filter := range filter.filters {
				if filter.FilterRepository(repository) == false {
					allow = false
					if filterdef, ok := filter.(FilterDefinition); ok {
						log.Printf("Skipping repository \"%s\" (filtered by %v)", name, filterdef.Name())
					}
					break
				}
			}

			if allow {
				log.Printf("Found repository \"%s\"", name)
				no_of_repositories++
				reposChannel <- repository
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
