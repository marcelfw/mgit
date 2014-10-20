package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/marcelfw/mgit/cmd_status"
	"github.com/marcelfw/mgit/cmd_pwd"
	"github.com/marcelfw/mgit/repository"
	go_ini "github.com/vaughan0/go-ini"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// number of parallel processors.
const numDigesters = 5

// repositoryFilter defines a filter for repositories.
type repositoryFilter struct {
	rootDirectory string

	branch string
	remote string
}

// command is the interface used for each command.
type command interface {
	Usage(int) string
	Help() string

	Init(args []string)

	Run(repository.Repository) repository.Repository

	OutputHeader() []string
	Output(repository.Repository) []string
}

// analysePath extracts repositories from regular file paths.
func analysePath(filter repositoryFilter, reposChannel chan repository.Repository) filepath.WalkFunc {
	no_of_repositories := 0

	return func(vpath string, f os.FileInfo, err error) error {
		base := path.Base(vpath)
		if base == ".git" {
			// Name is Git-directory without rootDirectory.
			name := strings.TrimLeft(path.Dir(vpath)[len(filter.rootDirectory):], "/")
			if repository, ok := repository.NewRepository(no_of_repositories, name, vpath); ok {
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
func findRepositories(filter repositoryFilter) chan repository.Repository {
	reposChannel := make(chan repository.Repository, numDigesters)

	go func() {
		filepath.Walk(filter.rootDirectory, analysePath(filter, reposChannel))

		close(reposChannel)
	}()

	return reposChannel
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

// readShortcutFromConfiguration reads the configuration and return the filter for the shortcut.
// return bool false if something went wrong.
func readShortcutFromConfiguration(shortcut string, filter repositoryFilter) (repositoryFilter, bool) {
	user, err := user.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, "Cannot determine home directory!")
		return filter, false
	}

	filename := user.HomeDir + "/.mgit"
	if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
		config, err := go_ini.LoadFile(filename)
		if err != nil {
			fmt.Fprint(os.Stderr, "Cannot read configuration file, incorrect format!\n")
			return filter, false
		}

		r, _ := regexp.Compile("shortcut \"(.+)\"")
		for name, vars := range config {
			match := r.FindStringSubmatch(name)
			if len(match) >= 2 && match[1] == shortcut {
				for key, value := range vars {
					lkey := strings.ToLower(key)
					switch {
					case lkey == "rootdirectory":
						filter.rootDirectory = value
					case lkey == "remote":
						filter.remote = value
					case lkey == "branch":
						filter.branch = value
					}
				}

				return filter, true
			}
		}
	} else {
		fmt.Fprintf(os.Stderr, "Cannot find configuration file, looked for %s! (%v %v)\n", filename, fi, err)
		return filter, false
	}

	fmt.Fprintf(os.Stderr, "Could not find shortcut \"%s\"!\n", shortcut)
	return filter, false
}

// getCommands fetches all commands available for this run.
func getCommands() (commands map[string]command) {
	commands = make(map[string]command)

	commands["status"] = cmd_status.NewCommand()
	commands["pwd"] = cmd_pwd.NewCommand()

	return
}

// parseCommandline parses and validates the command-line and returns useful structs to continue.
func parseCommandline() (command string, args []string, filter repositoryFilter, ok bool) {
	var rootDirectory string
	var remote string
	var branch string
	var shortcut string

	filter = repositoryFilter{rootDirectory: "."}
	flag.StringVar(&rootDirectory, "root", "", "set root directory")
	flag.StringVar(&remote, "r", "", "set remote to filter")
	flag.StringVar(&branch, "b", "", "set branch to filter")
	flag.StringVar(&shortcut, "s", "", "read settings with name from configuration file")
	flag.Parse()

	if flag.NArg() == 0 {
		return command, args, filter, false
	}

	if shortcut != "" {
		filter, ok = readShortcutFromConfiguration(shortcut, filter)
		if !ok {
			return command, args, filter, false
		}
	}

	if rootDirectory != "" {
		filter.rootDirectory = rootDirectory
	}
	if remote != "" {
		filter.remote = remote
	}
	if branch != "" {
		filter.branch = branch
	}

	args = flag.Args()
	command = args[0]
	args = args[1:]

	return command, args, filter, true
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
func runCommand(command command, filter repositoryFilter) {
	// Find repositories which match filter and put on inchannel.
	inChannel := findRepositories(filter)

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

	text_command, args, filter, ok := parseCommandline()
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
