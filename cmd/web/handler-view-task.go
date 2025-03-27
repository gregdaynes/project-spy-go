package main

import (
	"net/http"
)

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

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	data := app.newTemplateData(r)

	data.SearchData = app.searchData()
	data.TaskLanes = app.renderTaskLanes()

	task := app.taskLanes[lane].Tasks[filename]

	data.CurrentTask = ViewTaskModel{
		Title:          task.Title,
		Body:           task.RawContents,
		ShowDetails:    task.HasPriorityOrTags(),
		Priority:       task.Priority,
		Tags:           task.Tags,
		AvailableLanes: getAvailableLanes(&task, app.taskLanes),
		Actions:        getAvailableActions(&task),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
