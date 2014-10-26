// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package repository implements detection, filtering and structure of repositories.
// This source defines various interfaces.
package repository

import "flag"


// FilterDefinition is the interface used for each filter definition.
type FilterDefinition interface {
	Usage() string // short string describing the usage

	// Add flags for the command-line parser.
	AddFlags(*flag.FlagSet) (Filter)
}

// Filter is the actual interface of a repository filter.
type Filter interface {
	Dump() string

// Return true if the repository should be included.
	FilterRepository(Repository) (bool)
}

// Command is the interface used for each command.
type Command interface {
	Usage() string // short string describing the usage
	Help() string  // help info

	Init(args []string) (Command)

	// Return true if command can be interactive.
	IsInteractive() bool

	// Run the actual command.
	Run(Repository) (Repository, bool)

	OutputHeader() []string
	Output(Repository) interface{}
}
