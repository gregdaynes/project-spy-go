package main

import (
	"net/http"
)

func (app *application) newTask(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = app.searchData()
	data.TaskLanes = app.renderTaskLanes()

	task := Task{}

	data.CurrentTask = ViewTaskModel{
		Title:          "",
		Body:           "",
		ShowDetails:    true,
		Priority:       0,
		Tags:           []string{},
		AvailableLanes: getAvailableLanes(&task, app.taskLanes),
		Actions:        getAvailableActions(&task),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
