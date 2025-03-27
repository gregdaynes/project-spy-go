package main

import (
	"net/http"

	"projectspy.dev/internal/search"
	"projectspy.dev/internal/task"
)

func (app *application) info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	data.ShowInfo = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
