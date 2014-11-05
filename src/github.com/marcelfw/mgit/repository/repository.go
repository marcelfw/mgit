// Copyright (c) 2014 Marcel Wouters

// Package repository implements detection, filtering and structure of repositories.
// This source structures a single repository for processing by others.
package repository

import (
	"bytes"
	"fmt"
	go_ini "github.com/vaughan0/go-ini"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"
)

type Repository struct {
	index int    // order in which repository was found
	Name  string // assumed name of the repo

	path    string // root work directory
	gitRoot string // actual git location

	currentBranch string // store the current branch
	status        string // store the porcelain status

	remotes  map[string]string // remotes
	branches []string          // branch names

	info map[string]interface{} // let commands store info from a run here
}

type ByIndex []Repository

func (a ByIndex) Len() int           { return len(a) }
func (a ByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByIndex) Less(i, j int) bool { return a[i].index < a[j].index }

var remoteRegexp *regexp.Regexp

// init
func init() {
	remoteRegexp = regexp.MustCompile("remote \"(.+)\"")
}

// NewRepository returns a Repository structure describing the repository.
// If there is an error, ok will be false.
func NewRepository(index int, name, gitpath string) (repository Repository, ok bool) {
	ok = false
	repository.index = index
	repository.Name = name
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

		repository.updateRemotes()
		repository.updateBranches()
	}
	return
}

// findRemotes fills the remotes array with all the names (of the remotes).
func (repository *Repository) updateRemotes() {
	remotes := make(map[string]string)

	if fi, err := os.Stat(repository.gitRoot + "/config"); err == nil && !fi.IsDir() {
		config, err := go_ini.LoadFile(repository.gitRoot + "/config")
		if err == nil {
			for name, vars := range config {
				match := remoteRegexp.FindStringSubmatch(name)
				if len(match) >= 2 {
					name := match[1]
					if value, ok := vars["url"]; ok {
						remotes[name] = value
					}
				}
			}
		}

		repository.remotes = remotes
	}
}

// findBranches fills the branches array with all the names (of the branches).
func (repository *Repository) updateBranches() {
	branches := make([]string, 0, 10)

	if fi, err := os.Stat(repository.gitRoot + "/logs/refs/heads"); err == nil && fi.IsDir() {
		if fis, err := ioutil.ReadDir(repository.gitRoot + "/logs/refs/heads"); err == nil {
			for _, fi := range fis {
				// We don't support branches in subdirectories.
				if !fi.IsDir() {
					branches = append(branches, fi.Name())
				}
			}
		}

		repository.branches = branches
	} else {
		fmt.Printf("! no directory [%v]", err)
	}
}

func (repository Repository) ExecGit(args ...string) (result string, ok bool) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repository.path

	output, err := cmd.CombinedOutput()
	if err == nil {
		return string(output), true
	}

	return "", false
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
	if branch, ok := repository.ExecGit("rev-parse", "--abbrev-ref", "HEAD"); ok {
		repository.currentBranch = strings.TrimRight(branch, "\r\n")
	}
	repository.status, _ = repository.ExecGit("status", "--porcelain")

}

// fetchRemote performs a fetch of a specific remote.
func (repository *Repository) fetchRemote(remote string) {
	_, _ = repository.ExecGit("fetch", remote)
}

// PathMatch returns true if path matches.
func (repository *Repository) PathMatch(match string) bool {
	if strings.Index(repository.path, match) >= 0 {
		return true
	}
	return false
}

// IsBranch return true if branch is a branch.
func (repository *Repository) IsBranch(branch string) bool {
	for _, b := range repository.branches {
		if b == branch {
			return true
		}
	}
	return false
}

// IsRemote returns true if remote is a remote.
func (repository *Repository) IsRemote(remote string) bool {
	_, ok := repository.remotes[remote]
	return ok
}

// RemotePathContains returns true if remote path contains search.
func (repository *Repository) RemotePathContains(search string) bool {
	for _, r := range repository.remotes {
		if strings.Contains(r, search) {
			return true
		}
	}
	return false
}

// NameContains returns true if name contains search.
func (repository *Repository) NameContains(search string) bool {
	if strings.Contains(repository.Name, search) {
		return true
	}
	return false
}

// GetPath returns repository root directory.
func (repository *Repository) GetPath() string {
	return repository.path
}

// GetCurrentBranch returns the current branch.
func (repository *Repository) GetCurrentBranch() string {
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
	macros["Name"] = repository.Name
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
