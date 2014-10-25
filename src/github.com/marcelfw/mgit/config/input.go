// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config implements configuration and start-up.
// This source parses the command-line and reads additional input configuration.
package config

import (
	"flag"
	"github.com/marcelfw/mgit/repository"
	go_ini "github.com/vaughan0/go-ini"
	"os/user"
	"regexp"
	"fmt"
	"os"
	"strings"
	"strconv"
)

// readShortcutFromConfiguration reads the configuration and return the filter for the shortcut.
// return bool false if something went wrong.
func readShortcutFromConfiguration(shortcut string, filterMap map[string]string) (map[string]string, bool) {
	//filterMap = make(map[string]string)

	user, err := user.Current()
	if err != nil {
		fmt.Fprint(os.Stderr, "Cannot determine home directory!")
		return filterMap, false
	}

	filename := user.HomeDir + "/.mgit"
	if fi, err := os.Stat(filename); err == nil && !fi.IsDir() {
		config, err := go_ini.LoadFile(filename)
		if err != nil {
			fmt.Fprint(os.Stderr, "Cannot read configuration file, incorrect format!\n")
			return filterMap, false
		}


		r, _ := regexp.Compile("shortcut \"(.+)\"")
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
	} else {
		fmt.Fprintf(os.Stderr, "Cannot find configuration file, looked for %s! (%v %v)\n", filename, fi, err)
		return filterMap, false
	}

	fmt.Fprintf(os.Stderr, "Could not find shortcut \"%s\"!\n", shortcut)
	return filterMap, false
}

// parseCommandline parses and validates the command-line and return useful structs to continue.
func ParseCommandline() (command string, args []string, filter repository.RepositoryFilter, ok bool) {
	var rootDirectory string
	var depth int
	var remote string
	var noremote string
	var branch string
	var nobranch string
	var shortcut string

	filter = repository.RepositoryFilter{}

	preCommandFlags := flag.NewFlagSet("precommandflags", flag.ContinueOnError)
	preCommandFlags.StringVar(&rootDirectory, "root", "", "set root directory")
	preCommandFlags.IntVar(&depth, "d", 0, "maximum depth to search in")
	preCommandFlags.StringVar(&remote, "r", "", "select only with this remote")
	preCommandFlags.StringVar(&noremote, "nr", "", "select only without this remote")
	preCommandFlags.StringVar(&branch, "b", "", "select only with this branch")
	preCommandFlags.StringVar(&branch, "nb", "", "select only without this branch")
	preCommandFlags.StringVar(&shortcut, "s", "", "read settings with name from configuration file")

	preCommandFlags.Parse(os.Args[1:])

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