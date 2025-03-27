package main

import (
	"net/http"
)

type SearchEntry []string
type SearchData []SearchEntry

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = app.searchData()
	data.TaskLanes = app.renderTaskLanes()

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
