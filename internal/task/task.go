package task

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
	Actions         map[string]Action
	ShowDetails     bool
	Body            string
	AvailableLanes  map[string]TaskLane
}

type Action struct {
	Label  string
	Name   string
	Method string
	Action string
}

func (t *Task) HasPriorityOrTags() bool {
	if t.Priority > 0 || len(t.Tags) > 0 {
		return true
	}

	return false
}

func GetAvailableLanes(t *Task, lanes map[string]TaskLane) map[string]TaskLane {
	taskLanes := make(map[string]TaskLane)

	for name, lane := range lanes {
		taskLanes[name] = TaskLane{
			Name:     lane.Name,
			Slug:     lane.Slug,
			Selected: t.Lane == lane.Name,
		}
	}

	return taskLanes
}

func GetAvailableActions(t *Task, mode string) map[string]Action {
	actions := make(map[string]Action)

	switch mode {
	case "create":
		actions["save"] = Action{
			Label:  "Create",
			Name:   "save",
			Action: "/new/",
			Method: http.MethodPost,
		}
	case "view":
		actions["view"] = Action{
			Label:  "View",
			Name:   "view",
			Method: http.MethodGet,
			Action: "/view/" + t.RelativePath,
		}
	case "edit":
		actions["save"] = Action{
			Label:  "Update",
			Name:   "update",
			Action: "/update/" + t.RelativePath,
			Method: http.MethodPost,
		}
		actions["archive"] = Action{
			Label:  "Archive",
			Name:   "archive",
			Action: "/archive/" + t.RelativePath,
			Method: http.MethodGet,
		}
		actions["delete"] = Action{
			Label:  "Delete",
			Name:   "delete",
			Action: "/delete/" + t.RelativePath,
			Method: http.MethodGet,
		}
	}

	return actions
}
