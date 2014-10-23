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

	Init(args []string)

	Run(repository.Repository) repository.Repository

	OutputHeader() []string
	Output(repository.Repository) []string
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

				outChannel <- command.Run(repository)
			}
			wg.Done()
		}()
	}
	wg.Wait()
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
	cmds["pwd"] = commands.NewPwdCommand()

	return
}


// Output an text string table.
func outputTextTable(header []string, rows [][]string) string {
	var buffer bytes.Buffer

	// Storage for column widths and line.
	var column_width []int = make([]int, len(header))
	var line_columns []string = make([]string, len(header))

	// Init column width header columns.
	for idx, column := range header {
		column_width[idx] = len(column)
	}

	// Determine column widths.
	for _, row := range rows {
		for idx, column := range row {
			if len(column) > column_width[idx] {
				column_width[idx] = len(column)
			}
		}
	}

	// Fill line columns.
	for idx, _ := range header {
		line_columns[idx] = strings.Repeat("-", column_width[idx])
	}

	// Inserts header and lines into rows.
	rows = append(rows, header, header)
	copy(rows[2:], rows[0:len(rows)-1])
	rows[0] = header
	rows[1] = line_columns

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

	// Clear counter.
	fmt.Printf("\r        \r")

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

	var command command
	if command, ok = commands[text_command]; ok == false {
		Usage(commands)
		return
	}

	// Let the command initialize itself with the arguments.
	command.Init(args)

	// Run the actual command.
	runCommand(command, filter)
}
