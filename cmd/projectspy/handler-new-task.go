package main

import (
	"net/http"

	"projectspy.dev/internal/search"
	"projectspy.dev/internal/task"
	"projectspy.dev/web"
)

func (app *application) newTask(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	newTask := task.Task{}

	data.CurrentTask = web.ViewTaskModel{
		Title:          "",
		Body:           "",
		ShowDetails:    true,
		Priority:       0,
		Tags:           []string{},
		AvailableLanes: task.GetAvailableLanes(&newTask, app.taskLanes),
		Actions:        task.GetAvailableActions(&newTask),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
