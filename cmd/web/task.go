package main

import (
	"net/http"
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

func getAvailableLanes(t *Task, lanes map[string]TaskLane) map[string]ViewLaneModel {
	taskLanes := make(map[string]ViewLaneModel)

	for name, lane := range lanes {
		taskLanes[name] = ViewLaneModel{
			Name:     lane.Name,
			Slug:     lane.Slug,
			Selected: t.Lane == lane.Name,
		}
	}

	return taskLanes
}

func getAvailableActions(t *Task) map[string]ViewActionModel {
	actions := make(map[string]ViewActionModel)

	if t.Title == "" {
		actions["save"] = ViewActionModel{
			Label:  "Create",
			Name:   "save",
			Action: "/new/",
			Method: http.MethodPost,
		}

		return actions
	}

	actions["view"] = ViewActionModel{
		Label:  "View",
		Name:   "view",
		Method: http.MethodGet,
		Action: "/view/" + t.Lane + "/" + t.Filename,
	}

	actions["save"] = ViewActionModel{
		Label:  "Update",
		Name:   "update",
		Action: "/update/" + t.Lane + "/" + t.Filename,
		Method: http.MethodPost,
	}

	actions["archive"] = ViewActionModel{
		Label:  "Archive",
		Name:   "archive",
		Action: "/archive/" + t.Lane + "/" + t.Filename,
		Method: http.MethodGet,
	}
	actions["delete"] = ViewActionModel{
		Label:  "Delete",
		Name:   "delete",
		Action: "/delete/" + t.Lane + "/" + t.Filename,
		Method: http.MethodGet,
	}

	return actions
}
