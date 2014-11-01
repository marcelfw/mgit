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
	"strings"
)

type cmdEcho struct {
	args []string

	repository repository.Repository
}

func NewEchoCommand() cmdEcho {
	var cmd cmdEcho

	return cmd
}

func (cmd cmdEcho) Usage() string {
	return "Echo output after conversion."
}

func (cmd cmdEcho) Help() string {
	return `Echo output after conversion.

Parses the rest of the command-line and performs all macro conversion.
Useful for testing macros.`
}

func (cmd cmdEcho) Init(args []string) (outCmd repository.Command) {
	cmd.args = args
	return cmd
}

func (cmd cmdEcho) IsInteractive() bool {
	return false
}

func (cmd cmdEcho) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	return repository, true
}

func (cmd cmdEcho) OutputHeader() []string {
	return nil
}

func (cmd cmdEcho) Output(repository repository.Repository) interface{} {
	return strings.Join(repository.ReplaceMacros(cmd.args), " ")
}
