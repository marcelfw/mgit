// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns the last matching root directory of a repository search string.
package command

import (
	"github.com/marcelfw/mgit/engine"
	"github.com/marcelfw/mgit/repository"
	"os"
	"os/exec"
	"strings"
)

type cmdExec struct {
	args []string

	interactive bool
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

func (cmd cmdExec) Init(args []string, interactive bool) (outCmd repository.Command) {
	cmd.args = args
	cmd.interactive = interactive
	return cmd
}

func (cmd cmdExec) IsInteractive() bool {
	return cmd.interactive
}

func (cmd cmdExec) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	args := repository.ReplaceMacros(cmd.args)
	extCmd := exec.Command(args[0], args[1:]...)
	extCmd.Dir = repository.GetPath()

	repository.PutInfo("exec", "")

	if cmd.interactive {
		extCmd.Stdin = os.Stdin
		extCmd.Stdout = os.Stdout
		extCmd.Stderr = os.Stderr

		if err := extCmd.Start(); err != nil {
			repository.PutInfo("exec", err)
			return repository, false
		}

		if err := extCmd.Wait(); err != nil {
			repository.PutInfo("exec", err)
			return repository, false
		}

		repository.PutInfo("exec", "Ok")
	} else {
		result, err := extCmd.CombinedOutput()
		if err == nil {
			repository.PutInfo("exec", strings.TrimSpace(string(result)))
		}
	}
	return repository, true
}

func (cmd cmdExec) Header() []string {
	columns := make([]string, 2, 2)

	columns[0] = "Repository"
	columns[1] = "Output"

	return columns
}

func (cmd cmdExec) Output(repository repository.Repository) interface{} {
	return engine.FormatRow(repository.GetShowName(), repository.GetInfo("exec").(string))
}
