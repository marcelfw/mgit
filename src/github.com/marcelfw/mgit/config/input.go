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
func readShortcutFromConfiguration(shortcut string) (filterMap map[string]string, ok bool) {
	configs := findOrderedConfigs()

	filterMap = make(map[string]string)

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
func ParseCommandline(filterDefs []repository.FilterDefinition) (command string, args []string, repositoryFilter repository.RepositoryFilter, ok bool) {
	var rootDirectory string
	var depth int
	var shortcut string


	mgitFlags := flag.NewFlagSet("mgitFlags", flag.ContinueOnError)

	// These are truly hard-coded for now.
	mgitFlags.StringVar(&shortcut, "s", "", "read settings with name from configuration file")
	mgitFlags.StringVar(&rootDirectory, "root", "", "set root directory")
	mgitFlags.IntVar(&depth, "d", 0, "maximum depth to search in")

	filters := make([]repository.Filter, 0, len(filterDefs))
	for _, filterDef := range filterDefs {
		filters = append(filters, filterDef.AddFlags(mgitFlags))
	}

	mgitFlags.Parse(os.Args[1:])

	if mgitFlags.NArg() == 0 {
		return command, args, repositoryFilter, false
	}


	var filterMap map[string]string

	if shortcut != "" {
		filterMap, ok = readShortcutFromConfiguration(shortcut)
		if !ok {
			return command, args, repositoryFilter, false
		}
	}

	if rootDirectory == "" {
		if value, ok := filterMap["rootDirectory"]; ok {
			rootDirectory = value
		}
	}
	if depth == 0 {
		if value, ok := filterMap["depth"]; ok {
			if ivalue, err := strconv.ParseInt(value, 10, 0); err == nil {
				depth = int(ivalue)
			}
		}
	}

	mgitFlags.VisitAll(func(flag *flag.Flag){
		if value, ok := filterMap[flag.Name]; ok {
			fmt.Printf("There is a shortcut for [%v] '%s'\n", flag, value)
		}
	})

	repositoryFilter = repository.NewRepositoryFilter(rootDirectory, depth, filters)

	args = mgitFlags.Args()
	command = args[0]
	args = args[1:]

	return command, args, repositoryFilter, true
}

// createCommand creates a command based on a configuration section.
// returns _, false if command could not be created
func createCommand(vars map[string]string) (repository.Command, bool) {
	if value, ok := vars["git"]; ok {
		// add Git command
		return command.NewGitProxyCommand(value, vars), true
	}
	return nil, false
}

// AddConfigCommands add commands from the configuration files to the command list.
func AddConfigCommands(commands map[string]repository.Command) (map[string]repository.Command) {
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
