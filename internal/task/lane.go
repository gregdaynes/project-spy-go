package task

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"projectspy.dev/internal/config"
	event_bus "projectspy.dev/internal/event-bus"
)

type TaskLanes map[string]TaskLane

type TaskLane struct {
	Name     string
	Slug     string
	Path     string
	Tasks    Tasks
	Count    int
	Selected bool
}

func NewTaskLanes(config *config.Config) (TaskLanes, error) {
	var taskLanes = make(TaskLanes)

	files, err := os.ReadDir(".projectSpy")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, file := range files {
		for i, lane := range config.Lanes {
			if file.Name() == lane.Dir {
				taskLanes[file.Name()] = TaskLane{
					Slug: file.Name(),
					Name: file.Name(),
				}

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

		tasks := make(Tasks)

		for _, e := range entries {
			task, err := ParseFile(".projectSpy/" + lane.Slug + "/" + e.Name())
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

				if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
					lane, ok := lanes[laneName]
					if !ok {
						log.Fatal("lane not found", laneName)
					}

					_, ok = lane.Tasks[filename]
					if ok {
						delete(lane.Tasks, filename)
						eventBus.Publish("delete", filename)
					}

					task, err := ParseFile(event.Name)
					if err != nil {
						log.Fatal(err)
					}

					lanes[laneName].Tasks[filename] = task
					eventBus.Publish("update", laneName+"/"+filename)
				}

				if event.Has(fsnotify.Rename) || event.Has(fsnotify.Remove) {
					delete(lanes[laneName].Tasks, filename)
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

		lane.Path = "test"
		lanes[k] = lane
	}

	return watcher
}

func RenderTaskLanes(config *config.Config, lanes map[string]TaskLane) map[int]TaskLane {
	taskLanes := make(map[int]TaskLane)
	configLanes := config.Lanes

	for i := 0; i < len(configLanes); i++ {
		configLane := configLanes[i]

		if configLane.HasDir == false {
			continue
		}

		dir := configLane.Dir
		name := configLane.Name
		lane := lanes[dir]

		for _, task := range lane.Tasks {
			task.Actions = GetAvailableActions(&task, "view")
			lane.Tasks[task.Filename] = task
		}

		lane.Name = name
		lane.Count = len(lane.Tasks)
		taskLanes[i] = lane

	}

	return taskLanes
}
