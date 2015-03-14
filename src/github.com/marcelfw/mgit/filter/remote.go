// Copyright (c) 2014 Marcel Wouters

// Package filter implements all internal filters.
// This code filters on the presence of a remote or partial remote paths
package filter

import (
	"flag"
	"regexp"
	"strings"

	"github.com/marcelfw/mgit/repository"
)

type filterRemote struct {
	name string

	remote   *string
	noremote *string

	remoteurl   *string
	noremoteurl *string
}

var remoteRegexp *regexp.Regexp

// init
func init() {
	remoteRegexp = regexp.MustCompile("remote \"(.+)\"")
}

// NewRemoteFilter returns a new filterRemote filter.
func NewRemoteFilter() filterRemote {
	filter := filterRemote{name: "remote"}

	return filter
}

func (filter filterRemote) Name() string {
	return filter.name
}

func (filter filterRemote) Usage() map[string]string {
	return map[string]string{
		"-remote <remote>":           "Match when <remote> is found.",
		"-noremote <remote>":         "Match only when <remote> is not found.",
		"-remoteurl <partial-url>":   "Match when text matched <remoteurl>.",
		"-noremoteurl <partial-url>": "Match only when text does not match <remoteurl>.",
	}
}

func (filter filterRemote) AddFlags(flags *flag.FlagSet) repository.Filter {
	filter.remote = flags.String("remote", "", "select only with this remote")
	filter.noremote = flags.String("noremote", "", "select only without this remote")

	filter.remoteurl = flags.String("remoteurl", "", "select only with this value is found in the remote path")
	filter.noremoteurl = flags.String("noremoteurl", "", "select only when this value is not found in the remote path")

	return filter
}

// getRemotes returns remotes and paths.
func getRemotes(repository repository.Repository) (remotes map[string]string) {
	remotes = make(map[string]string)

	for name, vars := range repository.GetConfig() {
		match := remoteRegexp.FindStringSubmatch(name)
		if len(match) >= 2 {
			name := match[1]
			if value, ok := vars["url"]; ok {
				remotes[name] = value
			}
		}
	}

	return remotes
}

func (filter filterRemote) FilterRepository(repos repository.Repository) bool {
	remotes := getRemotes(repos)

	if *filter.remote != "" {
		if _, ok := remotes[*filter.remote]; !ok {
			return false
		}
	}
	if *filter.noremote != "" {
		if _, ok := remotes[*filter.noremote]; ok {
			return false
		}
	}

	if *filter.remoteurl != "" {
		found := false
		for _, path := range remotes {
			if strings.Contains(path, *filter.remoteurl) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	if *filter.noremoteurl != "" {
		for _, path := range remotes {
			if strings.Contains(path, *filter.noremoteurl) {
				return false
			}
		}
	}

	return true
}
