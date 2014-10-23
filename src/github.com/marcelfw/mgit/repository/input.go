// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repository implements detection, filtering and structure of repositories.
// This source parses the command-line and reads additional input configuration.
package repository

import (
	"flag"
	go_ini "github.com/vaughan0/go-ini"
	"os/user"
	"regexp"
	"strconv"

	"fmt"
	"os"
	"strings"
)

// readShortcutFromConfiguration reads the configuration and return the filter for the shortcut.
// return bool false if something went wrong.
func readShortcutFromConfiguration(shortcut string, filter RepositoryFilter) (RepositoryFilter, bool) {
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
					case lkey == "depth":
						depth, err := strconv.ParseInt(value, 10, 0)
						if err != nil {
							filter.depth = int(depth)
						} else {
							filter.depth = 0
						}
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

// parseCommandline parses and validates the command-line and returns useful structs to continue.
func ParseCommandline() (command string, args []string, filter RepositoryFilter, ok bool) {
	var rootDirectory string
	var depth int
	var remote string
	var noremote string
	var branch string
	var nobranch string
	var shortcut string

	filter = RepositoryFilter{rootDirectory: "."}
	flag.StringVar(&rootDirectory, "root", "", "set root directory")
	flag.IntVar(&depth, "d", 0, "maximum depth to search in")
	flag.StringVar(&remote, "r", "", "set remote to filter to include")
	flag.StringVar(&noremote, "nr", "", "set remote to filter to exclude")
	flag.StringVar(&branch, "b", "", "set branch to filter to include")
	flag.StringVar(&branch, "nb", "", "set branch to filter to exclude")
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
	if depth != 0 {
		filter.depth = depth
	}
	if remote != "" {
		filter.remote = remote
	}
	if noremote != "" {
		filter.remote = noremote
		filter.noremote = true
	}
	if branch != "" {
		filter.branch = branch
	}
	if nobranch != "" {
		filter.branch = nobranch
		filter.nobranch = true
	}

	args = flag.Args()
	command = args[0]
	args = args[1:]

	return command, args, filter, true
}
