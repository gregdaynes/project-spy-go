package task

import (
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"projectspy.dev/internal/config"
	event_bus "projectspy.dev/internal/event-bus"
)

type TaskLanes []TaskLane

// TODO remove this and replace with an interface for tasklane
type TaskLane struct {
	Name     string
	Title    string
	Slug     string
	Tasks    []Task
	Count    int
	Selected bool
}

func NewTaskLanes(config *config.Config) (TaskLanes, error) {
	var taskLanes = []TaskLane{}

	files, err := os.ReadDir(".projectSpy")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, file := range files {
		for i, lane := range config.Lanes {
			if file.Name() == lane.Dir {
				taskLanes = append(taskLanes, TaskLane{
					Slug:  file.Name(),
					Name:  file.Name(),
					Title: config.Lanes[i].Name,
				})

				config.Lanes[i].HasDir = true
			}
		}
	}

	return taskLanes, nil
}

func ListTasks(lanes TaskLanes) {
	for k, lane := range lanes {
		entries, err := os.ReadDir(".projectSpy/" + lane.Slug)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			task, err := ParseFile(".projectSpy/" + lane.Slug + "/" + e.Name())
			if err != nil {
				log.Fatal(err)
			}
			// task.Lanes = &lanes

			lane.Tasks = append(lane.Tasks, task)
		}

		lane.Count = len(lane.Tasks)
		lanes[k] = lane
	}
}

func SetupWatcher(eventBus *event_bus.EventBus[string], lanes TaskLanes) *fsnotify.Watcher {
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
				name := strings.TrimPrefix(event.Name, ".projectSpy/")
				laneName := strings.Split(name, "/")[0]
				filename := strings.Split(name, "/")[1]

				i := slices.IndexFunc(lanes, func(lane TaskLane) bool {
					return lane.Slug == laneName
				})

				j := slices.IndexFunc(lanes[i].Tasks, func(task Task) bool {
					return task.Filename == filename
				})

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
					// delete element from slice
					if j != -1 {
						lanes[i].Tasks = append(lanes[i].Tasks[:j], lanes[i].Tasks[j+1:]...)
					}
					eventBus.Publish("delete", filename)

					task, err := ParseFile(event.Name)
					if err != nil {
						log.Fatal(err)
					}

					lanes[i].Tasks = append(lanes[i].Tasks, task)
					eventBus.Publish("update", laneName+"/"+filename)
				}

				if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
					if j != -1 {
						lanes[i].Tasks = append(lanes[i].Tasks[:j], lanes[i].Tasks[j+1:]...)
					}
					eventBus.Publish("remove", name)
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

		lanes[k] = lane
	}

	return watcher
}

func RenderTaskLanes(config *config.Config, lanes []TaskLane) map[int]TaskLane {
	taskLanes := make(map[int]TaskLane)
	configLanes := config.Lanes

	for i := 0; i < len(configLanes); i++ {
		configLane := configLanes[i]

		if configLane.HasDir == false {
			continue
		}

		j := slices.IndexFunc(lanes, func(lane TaskLane) bool {
			return lane.Slug == configLane.Dir
		})
		lane := lanes[j]

		dir := configLane.Dir
		name := configLane.Name

		newLane := TaskLane{
			Name:     name,
			Slug:     dir,
			Tasks:    make([]Task, 0),
			Selected: false,
		}

		for _, task := range lane.Tasks {
			newLane.Tasks = append(newLane.Tasks, task)
		}

		sort.SliceStable(newLane.Tasks, func(i, j int) bool {
			if newLane.Tasks[i].Priority == newLane.Tasks[j].Priority {
				return newLane.Tasks[j].ModifiedTime.Before(newLane.Tasks[i].ModifiedTime)
			}

			return newLane.Tasks[i].Priority > newLane.Tasks[j].Priority
		})

		newLane.Count = len(lane.Tasks)

		taskLanes[i] = newLane

	}

	return taskLanes
}
