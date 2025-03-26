package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gosimple/slug"
)

type SearchEntry []string
type SearchData []SearchEntry

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	searchData := SearchData{}
	for LaneName, Lane := range app.taskLanes {
		for fileName, Task := range Lane.Tasks {
			entry := SearchEntry{}
			entry = append(entry, strings.ToLower(Task.Title+" "+Task.Description))
			entry = append(entry, slug.Make(LaneName+"-"+fileName))
			searchData = append(searchData, entry)
		}
	}
	searchJSON, _ := json.Marshal(searchData)
	data.SearchData = string(searchJSON)

	lanes := app.config.Lanes
	for i := 0; i < len(lanes); i++ {
		dir := lanes[i].Dir
		lane := app.taskLanes[dir]

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
				ID:              task.ID,
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
