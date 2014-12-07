// Copyright (c) 2014 Marcel Wouters

// Package repository implements detection, filtering and structure of repositories.
// This source structures a single repository for processing by others.
package repository

import (
	"bytes"
	go_ini "github.com/vaughan0/go-ini"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

type Repository struct {
	index int    // order in which repository was found
	name  string // assumed name of the repo

	path    string // root work directory
	gitRoot string // actual git location

	currentBranch string // store the current branch
	status        string // store the porcelain status

	config go_ini.File // stored config

	info map[string]interface{} // let commands store info from a run here
}

type ByIndex []Repository

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].index < a[j].index }

// NewRepository returns a Repository structure describing the repository.
// If there is an error, ok will be false.
func NewRepository(index int, name, gitpath string) (repository Repository, ok bool) {
	ok = false
	repository.index = index
	repository.name = name
	repository.path = path.Dir(gitpath)

	if fi, err := os.Stat(gitpath); err == nil {
		switch {
		case fi.IsDir():
			repository.gitRoot = gitpath
		case !fi.IsDir() && (fi.Size() < 4096):
			if redirFile, err := ioutil.ReadFile(gitpath); err == nil {
				if bytes.IndexAny(redirFile, "gitdir: ") == 0 {
					repository.gitRoot = path.Clean(repository.path + "/" + strings.TrimRight(string(redirFile[8:]), "\r\n"))
				}
			}
		}
	}

	if repository.gitRoot != "" {
		ok = true

		repository.readConfig()
	}
	return
}

func (repository *Repository) readConfig() {
	if fi, err := os.Stat(repository.gitRoot + "/config"); err == nil && !fi.IsDir() {
		config, err := go_ini.LoadFile(repository.gitRoot + "/config")
		if err == nil {
			repository.config = config
		}
	}
}

func (repository *Repository) GetConfig() go_ini.File {
	return repository.config
}

func (repository Repository) ExecGit(args ...string) (result string, err error, ok bool) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repository.path

	log.Printf("[%s] executing git with arguments %v", repository.GetShowName(), args)

	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), nil, true
	}

	log.Printf("[%s] git exited with error %v \"%s\"", repository.GetShowName(), err, output)

	return string(output), err, false
}

func (repository Repository) ExecGitInteractive(args ...string) (ok bool) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repository.path
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Printf("Command.Start returned err: %v!", err)
		return false
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Command.Wait returned err: %v!", err)
		return false
	}

	return true
}

// retrieveBasics retrieves the current branch, status.
func (repository *Repository) RetrieveBasics() {
	if branch, _, ok := repository.ExecGit("rev-parse", "--abbrev-ref", "HEAD"); ok {
		repository.currentBranch = strings.TrimRight(branch, "\r\n")
	}
	repository.status, _, _ = repository.ExecGit("status", "--porcelain")

}

// NameContains returns true if name contains search.
func (repository *Repository) NameContains(search string) bool {
	if strings.Contains(repository.name, search) {
		return true
	}
	return false
}

// GetGitRoot returns repository .git root directory.
func (repository *Repository) GetGitRoot() string {
	return repository.gitRoot
}

// GetName returns repository name.
func (repository *Repository) GetShowName() string {
	if repository.name == "" {
		return "(root)"
	}
	return repository.name
}

// GetPath returns repository root directory.
func (repository *Repository) GetPath() string {
	return repository.path
}

// GetCurrentBranch returns the current branch.
func (repository *Repository) GetCurrentBranch() string {
	if repository.currentBranch == "" {
		repository.RetrieveBasics()
	}
	return repository.currentBranch
}

// GetStatusJudgement judges the current status.
// Basically just shows if we have staged, unstaged or untracked files.
func (repository *Repository) GetStatusJudgement() string {
	var staged bool
	var unstaged bool
	var untracked bool
	lines := strings.Split(repository.status, "\n")
	for _, line := range lines {
		if len(line) >= 2 {
			switch {
			case line[0] == '?' || line[1] == '?':
				untracked = true
			case line[0] != ' ':
				staged = true
			case line[1] != ' ':
				unstaged = true
			}
		}
	}
	judgements := make([]string, 0, 3)
	if staged {
		judgements = append(judgements, "Staged")
	}
	if unstaged {
		judgements = append(judgements, "Unstaged")
	}
	if untracked {
		judgements = append(judgements, "Untracked")
	}

	return strings.Join(judgements, ", ")
}

// PutInfo stores information a command wants to publish later.
func (repository *Repository) PutInfo(name string, value interface{}) {
	if repository.info == nil {
		repository.info = make(map[string]interface{})
	}
	repository.info[name] = value
}

// GetInfo retrieves information a command wants to publish.
func (repository *Repository) GetInfo(name string) interface{} {
	return repository.info[name]
}

// ReplaceMacros replaces macros from the arguments and returns the strings with replacements.
func (repository Repository) ReplaceMacros(args []string) (out []string) {
	out = make([]string, len(args))

	macros := make(map[string]string)
	macros["Name"] = repository.name
	macros["Path"] = repository.GetPath()
	macros["CurrentBranch"] = repository.GetCurrentBranch()

	for idx, arg := range args {
		out[idx] = ""
		if t, err := template.New("arg").Parse(arg); err == nil {
			b := new(bytes.Buffer)
			if err := t.Execute(b, macros); err == nil {
				out[idx] = b.String()
			} else {
				log.Fatal(err)
			}
		}
	}

	return out
}
