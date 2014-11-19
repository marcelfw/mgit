// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns the last matching root directory of a repository search string.
package command

import (
	"github.com/marcelfw/mgit/engine"
	"github.com/marcelfw/mgit/repository"
	"os/exec"
	"strings"
)

type cmdExec struct {
	args []string
}

func NewExecCommand() cmdExec {
	var cmd cmdExec

	return cmd
}

func (cmd cmdExec) Usage() string {
	return "Execute a command."
}

func (cmd cmdExec) Help() string {
	return `Execute a command.

Performs macro conversion and runs the command(s).`
}

func (cmd cmdExec) Init(args []string) (outCmd repository.Command) {
	cmd.args = args
	return cmd
}

func (cmd cmdExec) IsInteractive() bool {
	return false
}

func (cmd cmdExec) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	args := repository.ReplaceMacros(cmd.args)
	extCmd := exec.Command(args[0], args[1:]...)
	extCmd.Dir = repository.GetPath()

	result, err := extCmd.CombinedOutput()
	if err == nil {
		repository.PutInfo("exec", strings.TrimSpace(string(result)))
	} else {
		repository.PutInfo("exec", "")
	}
	return repository, true
}

func (cmd cmdExec) Header() []string {
	columns := make([]string, 2, 2)

	columns[0] = "Script"
	columns[1] = "Output"

	return columns
}

func (cmd cmdExec) Output(repository repository.Repository) interface{} {
	return engine.FormatRow(repository.Name, repository.GetInfo("exec").(string))
}
