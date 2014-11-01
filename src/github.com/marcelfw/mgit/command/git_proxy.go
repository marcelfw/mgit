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
// This source returns a custom git command
package command

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
	"os/exec"
)

type cmdGitProxy struct {
	interactive bool

	command string
	args []string

	usage string
	help string
}

func NewGitProxyCommand(command string, vars map[string]string) cmdGitProxy {
	var cmd cmdGitProxy

	cmd.command = command
	cmd.args = make([]string, 0, 10)
	cmd.args = append(cmd.args, command)

	cmd.usage = "Run \"git " + command + "\"."

	if value, ok := vars["usage"]; ok {
		cmd.usage = value
	}
	if value, ok := vars["help"]; ok {
		cmd.help = value
	}

	if value, ok := vars["interactive"]; ok {
		if value == "yes" || value == "1" || value == "true" {
			cmd.interactive = true
		}
	}

	return cmd
}

func (cmd cmdGitProxy) Usage() string {
	return cmd.usage
}

func (cmd cmdGitProxy) Help() string {
	if cmd.help == "" {
		extCmd := exec.Command("git", "help", cmd.command)

		output, err := extCmd.CombinedOutput()
		if err == nil {
			return string(output)
		}

		return "No help information available."
	}
	return cmd.help
}

func (cmd cmdGitProxy) Init(args []string) (outCmd repository.Command) {
	cmd.args = append(cmd.args, args...)
	return cmd
}

func (cmd cmdGitProxy) IsInteractive() bool {
	return cmd.interactive
}

func (cmd cmdGitProxy) ForceInteractive() {
	cmd.interactive = true
}

func (cmd cmdGitProxy) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	var ok bool
	if cmd.interactive {
		ok = repository.ExecGitInteractive(cmd.args...)

		repository.PutInfo("proxy."+cmd.command, "(interactive command ran)")
	} else {
		var result string

		result, ok = repository.ExecGit(cmd.args...)

		repository.PutInfo("proxy."+cmd.command, strings.TrimSpace(result))
	}

	return repository, ok
}

func (cmd cmdGitProxy) OutputHeader() []string {
	columns := make([]string, 2, 2)

	columns[0] = strings.Title(cmd.command)
	columns[1] = "Output"

	return columns
}

// Output returns the result of the command
func (cmd cmdGitProxy) Output(repository repository.Repository) interface{} {

	output := repository.GetInfo("proxy."+cmd.command).(string)
	lines := strings.Split(output, "\n")

	switch {
	case len(lines) == 0 || (len(lines) == 1 && output == ""):
		columns := make([]string, 2, 2)
		columns[0] = repository.Name
		columns[1] = "<>"
		return columns
	case len(lines) == 1:
		columns := make([]string, 2, 2)
		columns[0] = repository.Name
		columns[1] = output
		return columns
	default:
		rows := make([][]string, 0, len(lines))
		for idx, line := range lines {
			columns := make([]string, 2, 2)
			var pre string // pre is used to hopefully make it easier to see the lines belong together
			switch {
			case idx == 0:
				columns[0] = repository.Name
				pre = "   "
			case idx == len(lines) - 1:
				pre = "\\_ "
			default:
				pre = "|  "
			}
			columns[1] = pre + strings.TrimSpace(line)
			rows = append(rows, columns)
		}
		return rows
	}

	return nil
}
