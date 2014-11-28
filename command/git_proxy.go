// Copyright (c) 2014 Marcel Wouters

// Package command implements all internal commands.
// This source returns a custom git command
package command

import (
	"github.com/marcelfw/mgit/engine"
	"github.com/marcelfw/mgit/repository"
	"os/exec"
	"strings"
)

type cmdGitProxy struct {
	interactive bool

	command string
	args    []string

	usage string
	help  string
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

func (cmd cmdGitProxy) Init(args []string, interactive bool) (outCmd repository.Command) {
	cmd.args = append(cmd.args, args...)
	cmd.interactive = interactive
	return cmd
}

func (cmd cmdGitProxy) IsInteractive() bool {
	return cmd.interactive
}

func (cmd cmdGitProxy) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	var ok bool

	args := repository.ReplaceMacros(cmd.args)

	if cmd.interactive {
		ok = repository.ExecGitInteractive(args...)

		repository.PutInfo("proxy."+cmd.command, "(interactive command ran)")
	} else {
		var result string

		result, ok = repository.ExecGit(args...)

		repository.PutInfo("proxy."+cmd.command, strings.TrimSpace(result))
	}

	return repository, ok
}

func (cmd cmdGitProxy) Header() []string {
	columns := make([]string, 2, 2)

	columns[0] = strings.Title(cmd.command)
	columns[1] = "Output"

	return columns
}

// Output returns the result of the command
func (cmd cmdGitProxy) Output(repository repository.Repository) interface{} {
	return engine.FormatRow(repository.GetShowName(), repository.GetInfo("proxy."+cmd.command).(string))
}
