package cmd_pwd

import (
	"github.com/marcelfw/mgit/repository"
	"strings"
)

const name string = "pwd"

type cmdStatus struct {
	match string

	repository repository.Repository
}

func NewCommand() cmdStatus {
	var cmd cmdStatus

	return cmd
}

func (cmd cmdStatus) Usage(name_len int) string {
	return name + strings.Repeat(" ", name_len-len(name)) + " Return the last matching directory."
}

func (cmd cmdStatus) Help() string {
	return "Returns repository directory."
}

func (cmd cmdStatus) Init(args []string) {
	cmd.match = args[0]
}

func (cmd cmdStatus) Run(repository repository.Repository) repository.Repository {
	if repository.PathMatch(cmd.match) {
		cmd.repository = repository
	}
	return repository
}

func (cmd cmdStatus) OutputHeader() []string {
	return []string {
		"Name", "Branch", "Status",
	}
}

func (cmd cmdStatus) Output(repository repository.Repository) []string {
	columns := make([]string, 3, 3)

	columns[0] = cmd.repository.Name
	columns[1] = cmd.repository.GetCurrentBranch()
	columns[2] = cmd.repository.GetStatusJudgement()

	return columns
}
