// Copyright (c) 2014 Marcel Wouters

// Package config implements configuration and start-up.
// This source configures mgit.
package config

import (
	"github.com/marcelfw/mgit/command"
	"github.com/marcelfw/mgit/engine"
	"github.com/marcelfw/mgit/filter"
	"github.com/marcelfw/mgit/repository"
	"sort"
)

// git commands non-interactive we automatically pass-through
var gitPassThru = []string{"status", "fetch", "push", "pull", "log", "commit", "add", "remote", "branch", "archive", "tag"}

// Usage returns the usage for the program.
func Usage(filters []repository.FilterDefinition, commands map[string]repository.Command) string {
	textBas := `usage: mgit [<filters>] <command> [<args>]

`

	filTable := make([][]string, 0, len(filters)*2+5)
	filTable = append(filTable, []string{"  -s <shortcut>", "Read shortcut for filters."})
	filTable = append(filTable, []string{"  -root <directory>", "Root directory to search from."})
	filTable = append(filTable, []string{"  -depth <depth>", "Maximum depth to search in."})
	filTable = append(filTable, []string{"  -debug", "Show debug output."})
	filTable = append(filTable, []string{"  -i", "Assume command is interactive."})
	for _, filter := range filters {
		for flag, help := range filter.Usage() {
			filTable = append(filTable, []string{"  " + flag, help})
		}

	}
	textFil := "Filters are:\n" + engine.ReturnTextTable(nil, filTable)

	cmdTable := make([][]string, 0, len(commands))

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		usage := make([]string, 2, 2)

		usage[0] = "  " + name
		usage[1] = commands[name].Usage()

		cmdTable = append(cmdTable, usage)
	}

	textCmd := "\nCommands are:\n" + engine.ReturnTextTable(nil, cmdTable)

	return textBas + textFil + textCmd
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
	cmds["exec"] = command.NewExecCommand()
	cmds["list"] = command.NewListCommand()
	cmds["version"] = command.NewVersionCommand()

	for _, gitCommand := range gitPassThru {
		cmds[gitCommand] = command.NewGitProxyCommand(gitCommand, map[string]string{})
	}

	cmds = AddConfigCommands(cmds)

	return cmds
}
