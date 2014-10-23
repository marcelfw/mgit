// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main glues everything together :-)
package main

import (
	"bytes"
	"fmt"
	"github.com/marcelfw/mgit/commands"
	"github.com/marcelfw/mgit/repository"
	"os"
	"sort"
	"strings"
	"sync"
)

// number of parallel processors.
const numDigesters = 5

// command is the interface used for each command.
type command interface {
	Usage(int) string
	Help() string

	Init(args []string) interface{}

	Run(repository.Repository) (repository.Repository, bool)

	OutputHeader() []string
	Output(repository.Repository) []string
}

// Usage returns the usage for the program.
func Usage(commands map[string]command) {
	fmt.Fprintf(os.Stderr, `usage: mgit [-s <shortcut-name>] [-root <root-directory>] [-b <branch>] [-r <remote>]
			<command> [<args>]

Commands are:

`)

	var name_len int
	for name, _ := range commands {
		if len(name) > name_len {
			name_len = len(name)
		}
	}

	for _, command := range commands {
		fmt.Fprintln(os.Stderr, command.Usage(name_len))
	}
}

// getCommands fetches all commands available for this run.
func getCommands() (cmds map[string]command) {
	cmds = make(map[string]command)

	cmds["status"] = commands.NewStatusCommand()
	cmds["path"] = commands.NewPathCommand()

	return
}

// goRepositories concurrently performs some actions on each repository.
func goRepositories(inChannel chan repository.Repository, outChannel chan repository.Repository, command command) {
	var wg sync.WaitGroup
	wg.Add(numDigesters)
	for i := 0; i < numDigesters; i++ {
		go func() {
			for repository := range inChannel {
				// Always require this information.
				repository.RetrieveBasics()

				if outRepository, output := command.Run(repository); output == true {
					outChannel <- outRepository
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// Output an text string table.
func outputTextTable(header []string, rows [][]string) string {
	var buffer bytes.Buffer

	// Storage for column widths and line.
	var column_width []int
	var line_columns []string

	// Init column width header columns.
	if header != nil {
		column_width = make([]int, len(header))
		line_columns = make([]string, len(header))

		for idx, column := range header {
			column_width[idx] = len(column)
		}
	}

	// Determine column widths.
	for _, row := range rows {
		if len(column_width) == 0 {
			column_width = make([]int, len(row))
			line_columns = make([]string, len(row))
		}
		for idx, column := range row {
			if len(column) > column_width[idx] {
				column_width[idx] = len(column)
			}
		}
	}

	if header != nil {
		// Fill line columns.
		for idx, _ := range header {
			line_columns[idx] = strings.Repeat("-", column_width[idx])
		}

		// Inserts header and lines into rows.
		rows = append(rows, header, header)
		copy(rows[2:], rows[0:len(rows)-1])
		rows[0] = header
		rows[1] = line_columns
	}

	// Write actual columns.
	for _, row := range rows {
		for idx, column := range row {
			if idx > 0 {
				buffer.WriteByte(32)
			}

			buffer.WriteString(column)
			if len(column) < column_width[idx] {
				buffer.WriteString(strings.Repeat(" ", column_width[idx]-len(column)))
			}
		}

		buffer.WriteString("\n")
	}

	return buffer.String()
}

// Run the actual command with the filter.
func runCommand(command command, filter repository.RepositoryFilter) {
	// Find repositories which match filter and put on inchannel.
	inChannel := repository.FindRepositories(filter, numDigesters)

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

	// Simplify repository output to rows.
	rows := make([][]string, len(repositories), len(repositories))
	for row_idx, repository := range repositories {
		output := command.Output(repository)

		rows[row_idx] = output
	}

	// Output nicely.
	fmt.Print(outputTextTable(command.OutputHeader(), rows))
}

func main() {
	commands := getCommands()

	text_command, args, filter, ok := repository.ParseCommandline()
	if ok == false {
		Usage(commands)
		return
	}

	var curCommand command
	if curCommand, ok = commands[text_command]; ok == false {
		Usage(commands)
		return
	}

	// Let the command initialize itself with the arguments.
	initResult := curCommand.Init(args)
	if newCommand, ok := initResult.(command); ok == true {
		curCommand = newCommand
	}

	// Run the actual command.
	runCommand(curCommand, filter)
}
