package task

import (
	"net/http"
	"time"

	"github.com/gosimple/slug"
)

type Tasks map[string]Task

type Task struct {
	Lane            string
	Title           string
	RawContents     string
	DescriptionHTML string
	Description     string //  Search data is derived from this
	Priority        int
	Tags            []string
	RelativePath    string
	Filename        string
	ModifiedTime    time.Time
	CreatedTime     time.Time
	Order           int
	Actions         map[string]Action // TODO remove this, it can be derived from a method
	AvailableLanes  []TaskLane
}

type Action struct {
	Label  string
	Name   string
	Method string
	Action string
}

func (t *Task) ID() string {
	return slug.Make(t.Lane + "-" + t.Filename)
}

func (t *Task) GetAvailableLanes(lanes []TaskLane) (taskLanes []TaskLane) {
	for i := 0; i < len(lanes); i++ {
		taskLanes = append(taskLanes, TaskLane{
			Title:    lanes[i].Title,
			Slug:     lanes[i].Slug,
			Selected: t.Lane == lanes[i].Slug,
		})
	}

	return taskLanes
}

func (t *Task) GetAvailableActions(mode string) map[string]Action {
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
