package main

import (
	"encoding/json"
	"strings"

	"github.com/gosimple/slug"
)

func (app *application) searchData() string {
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
	return string(searchJSON)
}
