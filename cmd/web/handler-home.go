package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	fmt.Println(app.taskLanes)
	lanes := app.config.Lanes
	for i := 0; i < len(lanes); i++ {
		dir := lanes[i].Dir
		lane := app.taskLanes[dir]

		fmt.Println(i, dir)

		data.TaskLanes[i] = ViewLaneModel{
			Name:  lanes[i].Name,
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

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
