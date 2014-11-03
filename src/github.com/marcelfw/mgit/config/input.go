// Copyright (c) 2014 Marcel Wouters

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
	"path"
	"regexp"
	"strconv"
	"strings"
)

var shortcutRegexp *regexp.Regexp
var commandRegexp *regexp.Regexp

// init
func init() {
	shortcutRegexp = regexp.MustCompile("shortcut \"(.+)\"")
	commandRegexp = regexp.MustCompile("command \"(.+)\"")
}

// getOrderedConfigFiles finds all configuration files and returns them in order.
func findOrderedConfigs() (configs []go_ini.File) {
	files := make([]string, 0, 10)

	// Follow parent directories and add all configurations.
	if wd, err := os.Getwd(); err == nil {
		for {
			filename := wd + "/.mgit"
			if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
				files = append(files, filename)
			}

			nwd := path.Dir(wd)
			if nwd == wd || nwd == "." {
				break
			}

			wd = nwd
		}
	}

	// Add configuration from user' directory.
	if user, err := user.Current(); err == nil {
		filename := user.HomeDir + "/.mgit"
		if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
			files = append(files, filename)
		}
	}

	configs = make([]go_ini.File, 0, 2)

	for _, filename := range files {
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

	for _, config := range configs {
		for name, vars := range config {
			match := shortcutRegexp.FindStringSubmatch(name)
			if len(match) >= 2 && match[1] == shortcut {
				for key, value := range vars {
					filterMap[strings.ToLower(key)] = value
				}

				return filterMap, true
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Could not find shortcut \"%s\"!\n", shortcut)
	return filterMap, false
}

// ParseCommandline parses and validates the command-line and return useful structs to continue.
func ParseCommandline(osArgs []string, filterDefs []repository.FilterDefinition) (command string, cmdInteractive bool, args []string, repositoryFilter repository.RepositoryFilter, ok bool) {
	var rootDirectory string
	var depth int
	var shortcut string
	var interactive bool

	mgitFlags := flag.NewFlagSet("mgitFlags", flag.ContinueOnError)

	// These are truly hard-coded for now.
	mgitFlags.StringVar(&shortcut, "s", "", "read settings with name from configuration file")
	mgitFlags.StringVar(&rootDirectory, "root", "", "set root directory")
	mgitFlags.IntVar(&depth, "depth", 0, "maximum depth to search in")
	mgitFlags.BoolVar(&interactive, "i", false, "run command interactively")

	filters := make([]repository.Filter, 0, len(filterDefs))
	for _, filterDef := range filterDefs {
		filters = append(filters, filterDef.AddFlags(mgitFlags))
	}

	mgitFlags.Parse(osArgs)

	var filterMap map[string]string

	if shortcut != "" {
		filterMap, ok = readShortcutFromConfiguration(shortcut)
		if !ok {
			return command, false, args, repositoryFilter, false
		}
	}

	if mgitFlags.NArg() == 0 {
		return command, false, args, repositoryFilter, false
	}

	mgitFlags.VisitAll(func(flag *flag.Flag) {
		if value, ok := filterMap[flag.Name]; ok {
			if flag.Value.String() == "" {
				flag.Value.Set(value)
			}
		}
	})

	if rootDirectory == "" {
		if value, ok := filterMap["root"]; ok {
			rootDirectory = value
		}
		if rootDirectory == "" {
			rootDirectory = "."
		}
	}
	if depth == 0 {
		if value, ok := filterMap["depth"]; ok {
			if ivalue, err := strconv.ParseInt(value, 10, 0); err == nil {
				depth = int(ivalue)
			}
		}
	}
	if interactive {
		cmdInteractive = true
	}

	repositoryFilter = repository.NewRepositoryFilter(rootDirectory, depth, filters)

	args = mgitFlags.Args()
	command = args[0]
	args = args[1:]

	return command, cmdInteractive, args, repositoryFilter, true
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
func AddConfigCommands(commands map[string]repository.Command) map[string]repository.Command {
	configs := findOrderedConfigs()

	for _, config := range configs {
		for name, vars := range config {
			match := commandRegexp.FindStringSubmatch(name)
			if len(match) >= 2 {
				if command, ok := createCommand(vars); ok {
					commands[match[1]] = command
				}
			}
		}
	}

	return commands
}
