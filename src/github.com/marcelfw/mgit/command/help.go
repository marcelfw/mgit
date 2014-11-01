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
// This source returns help.
package command

import (
	"github.com/marcelfw/mgit/repository"
)

type cmdHelp struct {
	IsHelp bool
}

func NewHelpCommand() cmdHelp {
	var cmd cmdHelp

	cmd.IsHelp = true

	return cmd
}

func (cmd cmdHelp) Usage() string {
	return "Show this help information."
}

func (cmd cmdHelp) Help() string {
	return `Show help information.

Add command as argument to help for more information about the command.`
}

func (cmd cmdHelp) Init(args []string) (outCmd repository.Command) {
	return nil
}

func (cmd cmdHelp) IsInteractive() bool {
	return false
}

func (cmd cmdHelp) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	return repository, false
}
