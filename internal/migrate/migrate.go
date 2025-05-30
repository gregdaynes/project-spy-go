package util

import (
	"log"
	"os"
	"strings"

	"projectspy.dev/internal/config"
	"projectspy.dev/internal/task"
	"projectspy.dev/internal/util"
)

func MigrateProject() {
	cfg, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	taskLanes, err := task.NewTaskLanes(&cfg)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}

	for _, lane := range taskLanes {
		entries, err := os.ReadDir(".projectSpy/" + lane.Slug)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			dir := ".projectSpy/" + lane.Slug + "/"
			path := dir + e.Name()

			// If a task file doesn't have an id :abcxyz
			// rename the file before parsing.
			// This is temporary code handle projects with existing tasks. Maybe it should be a command to run instead?
			if !strings.Contains(e.Name(), ":") {
				// generate new id
				uId := util.GenerateId(e.Name())
				fnParts := strings.Split(e.Name(), ".")
				newName := fnParts[0] + ":" + uId + "." + fnParts[1]
				newPath := dir + newName

				e := os.Rename(path, newPath)
				if e != nil {
					log.Fatal(e)
				}
			}
		}
	}

	log.Println("Migration complete")
}
