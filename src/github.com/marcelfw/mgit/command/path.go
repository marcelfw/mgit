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

// Package command implements all internal commands.
// This source returns the last matching root directory of a repository search string.
package command

import (
	"github.com/marcelfw/mgit/repository"
)

type cmdPath struct {
	match string

	repository repository.Repository
}

func NewPathCommand() cmdPath {
	var cmd cmdPath

	return cmd
}

func (cmd cmdPath) Usage() string {
	return "Return repository path of matched names."
}

func (cmd cmdPath) Help() string {
	return `Return repository path of matched names.

Returns all the actual paths matched by the search argument.`
}

func (cmd cmdPath) Init(args []string) (outCmd repository.Command) {
	if len(args) >= 1 {
		cmd.match = args[0]
	}
	return cmd
}

func (cmd cmdPath) IsInteractive() bool {
	return false
}

func (cmd cmdPath) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	if repository.PathMatch(cmd.match) {
		return repository, true
	}
	return repository, false
}

func (cmd cmdPath) OutputHeader() []string {
	return nil
}

func (cmd cmdPath) Output(repository repository.Repository) interface{} {
	return repository.GetPath()
}
