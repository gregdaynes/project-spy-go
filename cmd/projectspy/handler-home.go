package main

import (
	"net/http"

	"projectspy.dev/internal/search"
	"projectspy.dev/internal/task"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
