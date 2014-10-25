// Copyright 2014 Marcel Wouters. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package commands implements all internal commands.
// This source returns a custom git command
package commands

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
)

type cmdGitProxy struct {
	concurrent bool

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

	cmd.concurrent = true
	if value, ok := vars["concurrent"]; ok {
		if value == "no" || value == "0" || value == "false" {
			cmd.concurrent = false
		}
	}

	return cmd
}

func (cmd cmdGitProxy) Usage() string {
	return cmd.usage
}

func (cmd cmdGitProxy) Help() string {
	return cmd.help
}

func (cmd cmdGitProxy) Init(args []string) (outCmd interface{}) {
	cmd.args = append(cmd.args, args...)
	return cmd
}

func (cmd cmdGitProxy) RunConcurrently() bool {
	return cmd.concurrent
}

func (cmd cmdGitProxy) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	result, ok := repository.ExecGit(cmd.args...)

	repository.PutInfo("proxy."+cmd.command, strings.TrimSpace(result))

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
