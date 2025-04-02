package task

import (
	"net/http"
	"time"

	"github.com/gosimple/slug"
)

type Task struct {
	Lane         string
	Title        string
	RawContents  string
	Description  string
	Priority     int
	Tags         []string
	RelativePath string
	Filename     string
	ModifiedTime time.Time
	CreatedTime  time.Time
	Order        int
	Lanes        []Lane
}

type Action struct {
	Label  string
	Name   string
	Method string
	Action string
}

func (t Task) ID() string {
	return slug.Make(t.Lane + "-" + t.Filename)
}

func (t Task) AvailableLanes() (taskLanes []Lane) {
	for i := 0; i < len(t.Lanes); i++ {
		taskLanes = append(taskLanes, Lane{
			Title:    t.Lanes[i].Title,
			Slug:     t.Lanes[i].Slug,
			Selected: t.Lane == t.Lanes[i].Slug,
		})
	}

	return taskLanes
}

func (t Task) GetAvailableActions(mode string) map[string]Action {
	actions := make(map[string]Action)

	if t.Title == "" {
		mode = "create"
	}

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
