package util

import (
	"fmt"
	"log"
	"os"
	"regexp"
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

			// Migration: Add ID to task without id suffix ":ABCXYZ"
			if !strings.Contains(e.Name(), ":") {
				// If a task file doesn't have an id :abcxyz
				// rename the file before parsing.
				// This is temporary code handle projects with existing tasks. Maybe it should be a command to run instead?
				// generate new id
				uId := util.GenerateId(e.Name())
				fnParts := strings.Split(e.Name(), ".")
				newName := fnParts[0] + ":" + uId + "." + fnParts[1]
				newPath := dir + newName

				err := os.Rename(path, newPath)
				if err != nil {
					log.Fatal(err)
				}

				path = newPath

				fmt.Printf("renamed file %s -> %s", e.Name(), newName)
			}

			contents, err := os.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}

			// Migration: Update syntax for changelogs
			re := regexp.MustCompile(`(?m)---\n\n(^\d{4}-\d{2}-\d{2} \d{2}:\d{2}\t.*$[\n]?)+`)
			rePrefix := regexp.MustCompile(`(?m)^`)
			for _, match := range re.FindAll(contents, -1) {
				str := strings.TrimPrefix(string(match), "---\n\n")
				str = rePrefix.ReplaceAllString(str, ":")
				str = "changelog\n" + str

				contents = re.ReplaceAll(contents, []byte(str))

				fmt.Printf("updated changelog syntax for %s\n", e.Name())
			}

			os.WriteFile(path, contents, 0644)
		}
	}

	log.Println("Migration complete")
}
