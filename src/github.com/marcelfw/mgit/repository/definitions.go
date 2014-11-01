// Copyright (c) 2014 Marcel Wouters
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
// documentation files (the "Software"), to deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
// Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
// WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS
// OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT
// OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

	IsInteractive() bool // Return true if command can be interactive.

	// Run the actual command.
	Run(Repository) (Repository, bool)
}

// RowOutputCommand is a command which outputs rows.
type RowOutputCommand interface {
	OutputHeader() []string // Column headers.
	Output(Repository) interface{} //
}

// InteractiveCommand is a special command that could be run interactively.
type InteractiveCommand interface {
	ForceInteractive() // Force command to be run interactive
}

