// Copyright (c) 2014 Marcel Wouterså

// Package repository implements detection, filtering and structure of repositories.
// This source defines various interfaces.
package repository

import "flag"

// FilterDefinition is the interface used for each filter definition.
type FilterDefinition interface {
	Usage() string // short string describing the usage

	// Add flags for the command-line parser.
	AddFlags(*flag.FlagSet) Filter
}

// Filter is the actual interface of a repository filter.
type Filter interface {
	// Return true if the repository should be included.
	FilterRepository(Repository) bool
}

// Command is shared interface used for each command.
type Command interface {
	Usage() string // short string describing the usage
	Help() string  // help info

	Init(args []string) Command
}

// RepositoryCommand is the interface used commands that act on repositories.
type RepositoryCommand interface {
	IsInteractive() bool // Return true if command can be interactive.

	// Run the actual command.
	Run(Repository) (Repository, bool)
}

// RowOutputCommand is a command which outputs rows.
type RowOutputCommand interface {
	Header() []string // Column headers.

	Output(Repository) interface{} // [][]string, []string or string
}

// LineOutputCommand is a command which outputs lines.
type LineOutputCommand interface {
	Header() string
	Footer() string

	Output(Repository) string
}

// InteractiveCommand is a special command that could be run interactively.
type InteractiveCommand interface {
	ForceInteractive() // Force command to be run interactive
}

// InfoCommand is a special command that shows internal info.
type InfoCommand interface {
	Output(map[string]Command, string) string // commands and version string
}
