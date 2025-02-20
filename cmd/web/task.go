package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

type Tasks []Task

type Task struct {
	Name            string
	ID              string
	Title           string
	Lane            string
	DescriptionHTML string
	Tags            []string
	Priority        int
	Actions         []Action
	FullPath        string
	RelativePath    string
	Filename        string
	Description     string
	ModifiedTime    time.Time
	CreatedTime     time.Time
	Order           int
}

type Action struct {
	Label string
	Key   string
}

func listTasks(lanes TaskLanes) {
	for k, lane := range lanes {
		entries, err := os.ReadDir(".projectSpy/" + lane.Slug)
		if err != nil {
			log.Fatal(err)
		}

		tasks := []Task{}

		for _, e := range entries {
			fmt.Println(e.Name())

			task, err := parseFile(".projectSpy/" + lane.Slug + "/" + e.Name())
			if err != nil {
				log.Fatal(err)
			}

			task.Actions = append(task.Actions, Action{Label: "View", Key: "view"})

			tasks = append(tasks, task)
		}

		lane.Tasks = tasks
		lane.Count = len(lane.Tasks)
		lanes[k] = lane
	}
}
