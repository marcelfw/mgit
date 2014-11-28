// Copyright (c) 2014 Marcel Wouters

// Package engine implements the engine.
// This source is it (you know).
package engine

import (
	"fmt"
	"github.com/marcelfw/mgit/repository"
	"log"
	"sort"
	"sync"
)

// channel size for pushing repositories
const numCachedRepositories = 100

// number of parallel processors.
const numDigesters = 5

// goRepositories concurrently performs an action on each repository.
func goRepositories(inChannel chan repository.Repository, outChannel chan repository.Repository, command repository.RepositoryCommand) {
	digesters := numDigesters
	if command.IsInteractive() {
		digesters = 1
	}

	var wg sync.WaitGroup
	wg.Add(digesters)
	for i := 0; i < digesters; i++ {
		go func() {
			for repository := range inChannel {
				if outRepository, output := command.Run(repository); output == true {
					outChannel <- outRepository
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// Run the actual command with the filter.
func RunCommand(command repository.RepositoryCommand, filter repository.RepositoryFilter) {
	// Find repositories which match filter and put on inchannel.
	inChannel := repository.FindRepositories(filter, numCachedRepositories)

	// Get additional information about repositories and put on outChannel.
	outChannel := make(chan repository.Repository, numDigesters)
	go func() {
		goRepositories(inChannel, outChannel, command)
		close(outChannel)
	}()

	// Merge all repositories from the outChannel into slice.
	repositories := make([]repository.Repository, 0, 1000)
	for repository := range outChannel {
		repositories = append(repositories, repository)
	}

	// Sort repositories for logical output.
	sort.Sort(repository.ByIndex(repositories))

	// Repository output.
	if rowOutputCommand, ok := command.(repository.RowOutputCommand); ok {
		rows := make([][]string, 0, len(repositories))
		for _, repository := range repositories {
			output := rowOutputCommand.Output(repository)

			switch output.(type) {
			case string:
				rows = append(rows, []string{output.(string)})
			case []string:
				rows = append(rows, output.([]string))
			case [][]string:
				rows = append(rows, output.([][]string)...)
			default:
				log.Fatal("Unknown return type.")
			}

		}

		// Output nicely.
		fmt.Print(ReturnTextTable(rowOutputCommand.Header(), rows))
	} else if lineOutputCommand, ok := command.(repository.LineOutputCommand); ok {
		output := ""

		header := lineOutputCommand.Header()
		if header != "" {
			output += header + "\n"
		}
		for _, repository := range repositories {
			line := lineOutputCommand.Output(repository)
			if line != "" {
				output += line + "\n"
			}
		}
		footer := lineOutputCommand.Footer()
		if footer != "" {
			output += footer + "\n"
		}

		fmt.Print(output)
	}
}
