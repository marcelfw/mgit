// Copyright (c) 2014 Marcel Wouters

// Package filter implements all internal filters.
// This code filters on the presence of a tag.
package filter

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/marcelfw/mgit/repository"
)

type filterTag struct {
	name string

	tag   *string
	notag *string
}

// NewTagFilter returns a new filterTag filter.
func NewTagFilter() filterTag {
	filter := filterTag{name: "tag"}

	return filter
}

func (filter filterTag) Name() string {
	return filter.name
}

func (filter filterTag) Usage() map[string]string {
	return map[string]string{
		"-tag <tag>":   "Match when <tag> is found.",
		"-notag <tag>": "Match only when <tag> is not found.",
	}
}

func (filter filterTag) AddFlags(flags *flag.FlagSet) repository.Filter {
	filter.tag = flags.String("tag", "", "select only with this tag")
	filter.notag = flags.String("notag", "", "select only without this tag")

	return filter
}

// getTags returns the tags.
func getTags(repository repository.Repository) (tags map[string]bool) {
	tags = make(map[string]bool)

	if fi, err := os.Stat(repository.GetGitRoot() + "/refs/tags"); err == nil && fi.IsDir() {
		if fis, err := ioutil.ReadDir(repository.GetGitRoot() + "/refs/tags"); err == nil {
			for _, fi := range fis {
				// We don't support tags in subdirectories.
				if !fi.IsDir() {
					tags[fi.Name()] = true
				}
			}
		}
	} else {
		log.Printf("! no directory [%v]", err)
	}

	//log.Printf("Tags for repository \"%s\" => \"%v\"", repository.Name, tags)

	return tags
}

func (filter filterTag) FilterRepository(repos repository.Repository) bool {
	tags := getTags(repos)

	if *filter.tag != "" {
		if _, ok := tags[*filter.tag]; !ok {
			return false
		}
	}
	if *filter.notag != "" {
		if _, ok := tags[*filter.notag]; ok {
			return false
		}
	}

	return true
}
