package search

import (
	"encoding/json"
	"strings"

	"github.com/gosimple/slug"
	"projectspy.dev/internal/task"
)

type SearchEntry []string

func SearchData(taskLanes task.TaskLanes) string {
	searchData := make([]SearchEntry, 0)
	for LaneName, Lane := range taskLanes {
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
