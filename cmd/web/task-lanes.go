package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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

		tasks := make(Tasks)

		for _, e := range entries {
			//name := ".projectSpy/" + lane.Slug + "/" + e.Name()

			task, err := parseFile(".projectSpy/" + lane.Slug + "/" + e.Name())
			if err != nil {
				log.Fatal(err)
			}

			tasks[e.Name()] = task
		}

		lane.Tasks = tasks
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

				// TODO this is brittle.
				laneName := strings.Split(event.Name, "/")[1]

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
					fmt.Println("create or write", event.Name)
					lane, ok := lanes[laneName]
					if !ok {
						log.Fatal("lane not found", laneName)
					}

					_, ok = lane.Tasks[event.Name]
					if ok {
						fmt.Println("deleting", event.Name)
						delete(lane.Tasks, event.Name)
					}

					task, err := parseFile(event.Name)
					if err != nil {
						log.Fatal(err)
					}

					lanes[laneName].Tasks[event.Name] = task
				}

				if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
					fmt.Println("rename or remove", event.Name)
					delete(lanes[laneName].Tasks, event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	for k, lane := range lanes {
		err = watcher.Add(".projectSpy/" + lane.Name)
		if err != nil {
			log.Fatal(err)
		}

		lane.Path = "test"
		lanes[k] = lane
	}

	return watcher
}
