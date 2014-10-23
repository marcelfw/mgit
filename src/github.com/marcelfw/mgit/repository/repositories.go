package repository

import (
	"os"
	"fmt"
	"strings"
	"path/filepath"
	"path"
	"time"
)


// RepositoryFilter defines a filter for repositories.
type RepositoryFilter struct {
	rootDirectory string
	depth int

	branch string
	remote string
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
					found = repository.IsBranch(filter.branch)
				}
				if found == true && filter.remote != "" {
					found = repository.IsRemote(filter.remote)
				}

				if found {
					no_of_repositories++
					fmt.Printf("\r%c %d", "/-\\|"[time.Now().Second() % 4], no_of_repositories)
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
