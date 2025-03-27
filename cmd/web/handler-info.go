package main

import (
	"net/http"
)

func (app *application) info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = app.searchData()
	data.TaskLanes = app.renderTaskLanes()

	data.ShowInfo = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
