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
// This source returns the status of all repositories.
package command

import (
	"github.com/marcelfw/mgit/repository"
	"strconv"
	"strings"
	"time"
)

type cmdList struct {
}

func NewListCommand() cmdList {
	var cmd cmdList

	return cmd
}

func (cmd cmdList) Usage() string {
	return "List each repository with basic information."
}

func (cmd cmdList) Help() string {
	return `Return the status of each repository.

Shown are:
  Name     Shortened work directory of repository
  Branch   Current branch
  Status   Status summary of repository
  Commit   Last author commit date
  Subject  Subject of last commit`
}

func (cmd cmdList) Init(args []string) (outCmd repository.Command) {
	// we don't do anything
	return nil
}

// Return human readable time.
func (cmd cmdList) getHumanTime(atime time.Time) string {
	now := time.Now()
	diff := now.Sub(atime)

	switch {
	case diff.Hours() >= 6*24:
		return atime.Format("2006-01-02")
	case diff.Hours() >= 4 || now.Hour() < 4:
		return atime.Format("Monday, 15:04")
	}

	return atime.Format("Today, 15:04")
}

func (cmd cmdList) IsInteractive() bool {
	return false
}

func (cmd cmdList) Run(repository repository.Repository) (outRepository repository.Repository, output bool) {
	repository.RetrieveBasics()

	log, _ := repository.ExecGit("log", "--max-count=1", "--format=%an : %ae : %at : %s")
	results := strings.SplitN(strings.TrimRight(log, "\r\n"), " : ", 4)

	repository.PutInfo("list.name", "-")
	repository.PutInfo("list.email", "-")
	repository.PutInfo("list.time", "-")
	repository.PutInfo("list.subject", "-")

	if len(results) == 4 {
		repository.PutInfo("list.name", results[0])
		repository.PutInfo("list.email", results[1])
		if unixtime, err := strconv.ParseInt(results[2], 10, 0); err == nil {

			repository.PutInfo("list.time", cmd.getHumanTime(time.Unix(int64(unixtime), 0)))
			//repository.PutInfo("list.time", results[2])
		}
		repository.PutInfo("list.subject", results[3])
	}

	return repository, true
}

func (cmd cmdList) OutputHeader() []string {
	return []string{
		"Name", "Branch", "Status", "Last commit", "Subject",
	}
}

func (cmd cmdList) Output(repository repository.Repository) interface{} {
	columns := make([]string, 5, 5)

	columns[0] = repository.Name
	columns[1] = repository.GetCurrentBranch()
	columns[2] = repository.GetStatusJudgement()
	columns[3] = repository.GetInfo("list.time").(string)
	columns[4] = repository.GetInfo("list.subject").(string)

	return columns
}
