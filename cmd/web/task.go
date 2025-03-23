package main

import (
	"time"
)

type Tasks map[string]Task

type Task struct {
	Name            string
	ID              string
	Lane            string
	Title           string
	RawContents     string
	DescriptionHTML string
	Description     string
	Priority        int
	Tags            []string
	FullPath        string
	RelativePath    string
	Filename        string
	ModifiedTime    time.Time
	CreatedTime     time.Time
	Order           int
}

func (t *Task) HasPriorityOrTags() bool {
	if t.Priority > 0 || len(t.Tags) > 0 {
		return true
	}

	return false
}
