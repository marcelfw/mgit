package commands

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
)

const namePwd string = "pwd"

type cmdPwd struct {
	match string

	repository repository.Repository
}

func NewPwdCommand() cmdPwd {
	var cmd cmdPwd

	return cmd
}

func (cmd cmdPwd) Usage(name_len int) string {
	return namePwd + strings.Repeat(" ", name_len-len(namePwd)) + " Return the last matching directory."
}

func (cmd cmdPwd) Help() string {
	return "Returns repository directory."
}

func (cmd cmdPwd) Init(args []string) {
	cmd.match = args[0]
}

func (cmd cmdPwd) Run(repository repository.Repository) repository.Repository {
	if repository.PathMatch(cmd.match) {
		cmd.repository = repository
	}
	return repository
}

func (cmd cmdPwd) OutputHeader() []string {
	return []string {
		"Name", "Branch", "Status",
	}
}

func (cmd cmdPwd) Output(repository repository.Repository) []string {
	columns := make([]string, 3, 3)

	columns[0] = cmd.repository.Name
	columns[1] = cmd.repository.GetCurrentBranch()
	columns[2] = cmd.repository.GetStatusJudgement()

	return columns
}
