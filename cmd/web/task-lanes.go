package main

import (
	"log"
	"os"
)

type TaskLanes map[string]TaskLane

type TaskLane struct {
	Name string
}

func newTaskLanes() (TaskLanes, error) {
	var taskLanes = make(TaskLanes)

	files, err := os.ReadDir(".projectSpy")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, file := range files {
		taskLanes[file.Name()] = TaskLane{
			Name: file.Name(),
		}
	}

	return taskLanes, nil
}
