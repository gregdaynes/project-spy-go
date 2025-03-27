package main

import (
	"net/http"

	"projectspy.dev/internal/search"
	"projectspy.dev/internal/task"
	"projectspy.dev/web"
)

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	currentTask := app.taskLanes[lane].Tasks[filename]

	data.CurrentTask = web.ViewTaskModel{
		Title:          currentTask.Title,
		Body:           currentTask.RawContents,
		ShowDetails:    currentTask.HasPriorityOrTags(),
		Priority:       currentTask.Priority,
		Tags:           currentTask.Tags,
		AvailableLanes: task.GetAvailableLanes(&currentTask, app.taskLanes),
		Actions:        task.GetAvailableActions(&currentTask),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
