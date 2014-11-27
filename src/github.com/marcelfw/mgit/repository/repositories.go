// Copyright (c) 2014 Marcel Wouters

// Package repository implements detection, filtering and structure of repositories.
// This source detects and filters repositories.
package repository

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// RepositoryFilter defines a filter for repositories.
type RepositoryFilter struct {
	rootDirectory string
	depth         int

	filters []Filter
}

var regexpWorktree *regexp.Regexp

func init() {
	regexpWorktree = regexp.MustCompile("worktree = (.+)")
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

		if f.IsDir() {
			return nil
		}

		if base != "HEAD" {
			return nil
		}
		var configFileInfo os.FileInfo
		if configFileInfo, err = os.Stat(path.Dir(vpath) + "/config"); err != nil {
			return nil
		}
		if configFileInfo.Size() > 40960 {
			// ignore configs too big
			log.Printf("Ignoring repository with this huge configuration \"%s\" (%d bytes)", vpath, configFileInfo.Size())
			return nil
		}

		var content []byte
		if content, err = ioutil.ReadFile(path.Dir(vpath) + "/config"); err != nil {
			log.Printf("Could not read configuration \"%s\" (error %s)", vpath, err)
			return nil
		}

		// resolve submodule worktree
		match := regexpWorktree.FindStringSubmatch(string(content))
		if len(match) >= 2 {
			// we assume the submodule has a .git file here
			gitPath = path.Clean(gitPath + "/" + match[1] + "/.git")
			if fi, err := os.Stat(gitPath); err != nil {
				return nil
			} else {
				if fi.IsDir() {
					return nil
				}
			}
		}

		// Name is Git-directory without rootDirectory.
		if gitPath != ".git" {
			name = strings.TrimPrefix(gitPath, filter.rootDirectory)
			name = strings.TrimLeft(name, "/")
			name = strings.TrimSuffix(name, "/.git")
		}

		if filter.depth > 0 && (strings.Count(name, "/")+1) > filter.depth {
			// if depth limit is set, ignore directories too deep.
			log.Printf("Skipping repository \"%s\" (filtered by depth)", name)
			return nil
		}

		repository, foundRepository := NewRepository(no_of_repositories, name, gitPath)

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
