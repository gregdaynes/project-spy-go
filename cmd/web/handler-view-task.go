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

	lanes := app.config.Lanes
	for i := 0; i < len(lanes); i++ {
		dir := lanes[i].Dir
		lane := app.taskLanes[dir]

		data.TaskLanes[i] = ViewLaneModel{
			Name:  lane.Name,
			Slug:  lane.Slug,
			Tasks: make(map[string]ViewTaskModel),
			Count: len(lane.Tasks),
		}

		for _, task := range lane.Tasks {
			actions := make(map[string]ViewActionModel)
			actions["view"] = ViewActionModel{
				Label:  "View",
				Name:   "view",
				Action: "/view/" + task.Lane + "/" + task.Filename,
				Method: http.MethodGet,
			}

			data.TaskLanes[i].Tasks[task.Filename] = ViewTaskModel{
				Lane:            task.Lane,
				Title:           task.Title,
				DescriptionHTML: task.DescriptionHTML,
				Priority:        task.Priority,
				Tags:            task.Tags,
				Order:           task.Order,
				Actions:         actions,
			}
		}
	}

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
