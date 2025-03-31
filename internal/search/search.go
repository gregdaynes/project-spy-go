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

	for i := 0; i < len(taskLanes); i++ {
		lane := taskLanes[i]

		for fileName, Task := range lane.Tasks {
			entry := SearchEntry{}
			entry = append(entry, strings.ToLower(Task.Title+" "+Task.Description))
			entry = append(entry, slug.Make(lane.Name+"-"+fileName))
			searchData = append(searchData, entry)
		}
	}
	searchJSON, _ := json.Marshal(searchData)
	return string(searchJSON)
}
