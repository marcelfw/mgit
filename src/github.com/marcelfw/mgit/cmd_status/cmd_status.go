package cmd_status

import (
	"github.com/marcelfw/mgit/repository"
	"fmt"
	"strings"
)

const name string = "status"

type cmdStatus struct {
}

func NewCommand() (cmdStatus) {
	var cmd cmdStatus

	return cmd
}

func (cmd cmdStatus) Usage(name_len int) string {
	return name + strings.Repeat(" ", name_len - len(name)) + " Return the status of each repository"
}

func (cmd cmdStatus) Help() string {
	return "Returns really short status for repository."
}

func (cmd cmdStatus) Init(args []string) {
	// we don't do anything
}

func (cmd cmdStatus) Run(repository repository.Repository) (repository.Repository) {
	repository.RetrieveBasics()

	return repository
}

func (cmd cmdStatus) Output(repository repository.Repository) (string) {
	return fmt.Sprintf("%-30s %-20s  %-10s", repository.Name, repository.GetCurrentBranch(), repository.GetStatusJudgement())
}
