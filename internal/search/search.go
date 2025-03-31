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

		for j := 0; j < len(lane.Tasks); j++ {
			task := lane.Tasks[j]

			entry := SearchEntry{}
			entry = append(entry, strings.ToLower(task.Title+" "+task.Description))
			entry = append(entry, slug.Make(lane.Name+"-"+task.Filename))
			searchData = append(searchData, entry)
		}
	}
	searchJSON, _ := json.Marshal(searchData)
	return string(searchJSON)
}
