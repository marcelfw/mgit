// Copyright (c) 2014 Marcel Wouters

// Package config implements configuration and start-up.
// This source configures mgit.
package config

import (
	"fmt"
	"github.com/marcelfw/mgit/command"
	"github.com/marcelfw/mgit/engine"
	"github.com/marcelfw/mgit/filter"
	"github.com/marcelfw/mgit/repository"
	"os"
)

// git commands non-interactive we automatically pass-through
var gitPassThru = []string{"status", "fetch", "push", "pull", "log", "commit", "add", "remote", "branch", "archive", "tag"}

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

	fmt.Fprint(os.Stdout, engine.ReturnTextTable(nil, cmdTable))
}

// getFilters returns all filters.
func GetFilterDefs() []repository.FilterDefinition {
	filters := make([]repository.FilterDefinition, 0, 10)

	filters = append(filters, filter.NewNameFilter())
	filters = append(filters, filter.NewRemoteFilter())
	filters = append(filters, filter.NewBranchFilter())
	filters = append(filters, filter.NewTagFilter())

	return filters
}

// getCommands fetches all commands available for this run.
func GetCommands() map[string]repository.Command {
	cmds := make(map[string]repository.Command)

	cmds["help"] = command.NewHelpCommand()
	cmds["echo"] = command.NewEchoCommand()
	cmds["list"] = command.NewListCommand()
	cmds["version"] = command.NewVersionCommand()

	for _, gitCommand := range gitPassThru {
		cmds[gitCommand] = command.NewGitProxyCommand(gitCommand, map[string]string{})
	}

	cmds = AddConfigCommands(cmds)

	return cmds
}
