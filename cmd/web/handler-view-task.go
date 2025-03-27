package main

import (
	"net/http"
)

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
