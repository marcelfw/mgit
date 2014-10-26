// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config implements configuration and start-up.
// This source parses the command-line and reads additional input configuration.
package config

import (
	"flag"
	"fmt"
	"github.com/marcelfw/mgit/command"
	"github.com/marcelfw/mgit/repository"
	go_ini "github.com/vaughan0/go-ini"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
)

// Filter is the interface used for each filter.
type Filter interface {
	Usage() string // short string describing the usage

	// Add flags for the command-line parser.
	AddFlags(*flag.FlagSet)

	// Filter a single repository.
	// Return true if the repository should be included.
	FilterRepository(repository.Repository) (bool)
}

// Command is the interface used for each command.
type Command interface {
	Usage() string // short string describing the usage
	Help() string  // help info

	Init(args []string) interface{}

	// Return true if run can be executed concurrently.
	RunConcurrently() bool

	// Run the actual command.
	Run(repository.Repository) (repository.Repository, bool)

	OutputHeader() []string
	Output(repository.Repository) interface{}
}


// getOrderedConfigFiles finds all configuration files and returns them in order.
func findOrderedConfigs() (configs []go_ini.File) {
	configs = make([]go_ini.File, 0, 2)

	user, err := user.Current()
	if err != nil {
		// @todo panic or fail silently?
		fmt.Fprint(os.Stderr, "Cannot determine home directory!")
		return nil
	}

	filename := user.HomeDir + "/.mgit"
	if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
		config, err := go_ini.LoadFile(filename)
		if err != nil {
			// @todo panic or fail silently?
			fmt.Fprint(os.Stderr, "Cannot read configuration file, incorrect format!\n")
			return nil
		}

		configs = append(configs, config)
	}

	return configs
}

// readShortcutFromConfiguration reads the configuration and return the filter for the shortcut.
// return bool false if something went wrong.
func readShortcutFromConfiguration(shortcut string, filterMap map[string]string) (map[string]string, bool) {
	configs := findOrderedConfigs()

	r, _ := regexp.Compile("shortcut \"(.+)\"")
	for _, config := range configs {
		for name, vars := range config {
			match := r.FindStringSubmatch(name)
			if len(match) >= 2 && match[1] == shortcut {
				for key, value := range vars {
					lkey := strings.ToLower(key)
					switch {
					case lkey == "rootdirectory":
						filterMap["rootDirectory"] = value
					case lkey == "depth":
						filterMap["depth"] = value
					case lkey == "remote":
						filterMap["remote"] = value
					case lkey == "branch":
						filterMap["branch"] = value
					}
				}

				return filterMap, true
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Could not find shortcut \"%s\"!\n", shortcut)
	return filterMap, false
}

// ParseCommandline parses and validates the command-line and return useful structs to continue.
func ParseCommandline(filters []Filter) (command string, args []string, filter repository.RepositoryFilter, ok bool) {
	var rootDirectory string
	var depth int
	var remote string
	var noremote string
	var branch string
	var nobranch string
	var shortcut string

	preCommandFlags := flag.NewFlagSet("precommandflags", flag.ContinueOnError)
	preCommandFlags.StringVar(&rootDirectory, "root", "", "set root directory")
	preCommandFlags.IntVar(&depth, "d", 0, "maximum depth to search in")
	preCommandFlags.StringVar(&remote, "r", "", "select only with this remote")
	preCommandFlags.StringVar(&noremote, "nr", "", "select only without this remote")
	//preCommandFlags.StringVar(&branch, "b", "", "select only with this branch")
	//preCommandFlags.StringVar(&branch, "nb", "", "select only without this branch")
	preCommandFlags.StringVar(&shortcut, "s", "", "read settings with name from configuration file")

	for _, filter := range filters {
		filter.AddFlags(preCommandFlags)
	}

	preCommandFlags.Parse(os.Args[1:])

	for _, filter := range filters {
		fmt.Printf("filter[%v]\n", filter)
	}

	if preCommandFlags.NArg() == 0 {
		return command, args, filter, false
	}

	filterMap := make(map[string]string)

	if shortcut != "" {
		filterMap, ok = readShortcutFromConfiguration(shortcut, filterMap)
		if !ok {
			return command, args, filter, false
		}
	}

	if rootDirectory != "" {
		filterMap["rootDirectory"] = rootDirectory
	}
	if depth != 0 {
		filterMap["depth"] = strconv.FormatInt(int64(depth), 10)
	}
	if remote != "" {
		filterMap["remote"] = remote
	} else if noremote != "" {
		filterMap["noremote"] = noremote
	}
	if branch != "" {
		filterMap["branch"] = branch
	} else if nobranch != "" {
		filterMap["nobranch"] = branch
	}

	filter = repository.NewRepositoryFilter(filterMap)

	args = preCommandFlags.Args()
	command = args[0]
	args = args[1:]

	return command, args, filter, true
}

// createCommand creates a command based on a configuration section.
// returns _, false if command could not be created
func createCommand(vars map[string]string) (Command, bool) {
	if value, ok := vars["git"]; ok {
		// add Git command
		return command.NewGitProxyCommand(value, vars), true
	}
	return nil, false
}

// AddConfigCommands add commands from the configuration files to the command list.
func AddConfigCommands(commands map[string]Command) (map[string]Command) {
	configs := findOrderedConfigs()

	r, _ := regexp.Compile("command \"(.+)\"")
	for _, config := range configs {
		for name, vars := range config {
			match := r.FindStringSubmatch(name)
			if len(match) >= 2 {
				if command, ok := createCommand(vars); ok {
					commands[match[1]] = command
				}
			}
		}
	}

	return commands
}
