// Copyright (c) 2014 Marcel Wouters

// Package main glues everything together :-)
package main

import (
	"bytes"
	"fmt"
	"github.com/marcelfw/mgit/command"
	"github.com/marcelfw/mgit/config"
	"github.com/marcelfw/mgit/filter"
	"github.com/marcelfw/mgit/repository"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

// current version
var version = "0.0.1"

// channel size for pushing repositories
const numCachedRepositories = 100

// number of parallel processors.
const numDigesters = 5

// git commands non-interactive we automatically pass-through
var gitPassThru = []string{"status", "fetch", "push", "pull", "log", "commit", "add", "remote", "branch"}

// Usage returns the usage for the program.
func Usage(commands map[string]repository.Command) {
	fmt.Fprintf(os.Stderr, `usage: mgit [-s <shortcut-name>] [-root <root-directory>] -d <max-depth>
            [-branch <branch>] [-remote <remote>] [-nobranch <no-branch>] [-noremote <no-remote>]
            <command> [<args>]

Commands are:
`)

	cmdTable := make([][]string, 0, len(commands))

	for name, command := range commands {
		usage := make([]string, 2, 2)

		usage[0] = "  " + name
		usage[1] = command.Usage()

		cmdTable = append(cmdTable, usage)
	}

	fmt.Fprint(os.Stdout, returnTextTable(nil, cmdTable))
}

// getFilters returns all filters.
func getFilterDefs() []repository.FilterDefinition {
	filters := make([]repository.FilterDefinition, 0, 10)

	filters = append(filters, filter.NewBranchFilter())
	filters = append(filters, filter.NewRemoteFilter())
	filters = append(filters, filter.NewRemotePathFilter())
	filters = append(filters, filter.NewNameFilter())

	return filters
}

// getCommands fetches all commands available for this run.
func getCommands() map[string]repository.Command {
	cmds := make(map[string]repository.Command)

	cmds["help"] = command.NewHelpCommand()
	cmds["echo"] = command.NewEchoCommand()
	cmds["list"] = command.NewListCommand()
	cmds["version"] = command.NewVersionCommand()

	for _, gitCommand := range gitPassThru {
		cmds[gitCommand] = command.NewGitProxyCommand(gitCommand, map[string]string{})
	}

	cmds = config.AddConfigCommands(cmds)

	return cmds
}

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
				// Always require this information.
				//repository.RetrieveBasics()

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
func returnTextTable(header []string, rows [][]string) string {
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
				buffer.WriteString("  ")
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
func runCommand(command repository.RepositoryCommand, filter repository.RepositoryFilter) {
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
		fmt.Print(returnTextTable(rowOutputCommand.Header(), rows))
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

func main() {
	/*
		f, err := os.Create("mgit.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	*/

	filterDefs := getFilterDefs()
	commands := getCommands()

	textCommand, flagInteractive, args, filter, ok := config.ParseCommandline(os.Args[1:], filterDefs)
	if ok == false {
		Usage(commands)
		return
	}

	var curCommand repository.Command
	if curCommand, ok = commands[textCommand]; ok == false {
		Usage(commands)
		return
	}

	if flagInteractive {
		if iCommand, ok := curCommand.(repository.InteractiveCommand); ok {
			iCommand.ForceInteractive()
		}
	}

	// Let the command initialize itself with the arguments.
	initResult := curCommand.Init(args)
	// @todo remove this casting
	if newCommand, ok := initResult.(repository.Command); ok == true {
		curCommand = newCommand
	}

	/*
		if cmdType := reflect.TypeOf(curCommand); cmdType.Name() == "cmdHelp" {
			if len(args) == 1 {
				if helpCommand, ok := commands[args[0]]; ok == true {
					fmt.Fprintln(os.Stdout, helpCommand.Help())
					return
				}
			}
			Usage(commands)
			return
		}
	*/

	if repositoryCommand, ok := curCommand.(repository.RepositoryCommand); ok {
		// Run the actual command.
		runCommand(repositoryCommand, filter)
	} else if infoCommand, ok := curCommand.(repository.InfoCommand); ok {
		fmt.Fprintln(os.Stdout, infoCommand.Output(commands, version))
	} else {
		fmt.Fprintln(os.Stdout, curCommand.Help())
	}
}
