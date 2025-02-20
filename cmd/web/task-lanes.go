package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type TaskLanes map[string]TaskLane

type TaskLane struct {
	Name  string
	Slug  string
	Path  string
	Tasks Tasks
	Count int
}

type Tasks []Task

type Task struct {
	Name            string
	ID              string
	Title           string
	Lane            string
	DescriptionHTML string
	Tags            []string
	Priority        int
	Actions         []string
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
			Slug: file.Name(),
			Name: file.Name(),
		}
	}

	return taskLanes, nil
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

			tasks = append(tasks, Task{
				Name: e.Name(),
			})

			lane.Tasks = tasks
		}

		lane.Count = len(lane.Tasks)
		lanes[k] = lane
	}
}

func setupWatcher(lanes TaskLanes) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				fmt.Println("event:", event.Op, event.Name)

				// if event.Has(fsnotify.Write) {
				// 	log.Println("modified file:", event.Name)
				// }
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for k, lane := range lanes {
		// Add a path.
		err = watcher.Add(".projectSpy/" + lane.Name)
		if err != nil {
			log.Fatal(err)
		}

		lane.Path = "test"
		lanes[k] = lane
	}

	return watcher
}
