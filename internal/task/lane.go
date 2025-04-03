package task

import (
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/fsnotify/fsnotify"
	"projectspy.dev/internal/config"
	eventBus "projectspy.dev/internal/event-bus"
)

type Lanes []Lane

type Lane struct {
	Name     string
	Title    string
	Slug     string
	Tasks    []Task
	Count    int
	Selected bool
	Order    int
}

func NewTaskLanes(config *config.Config) (Lanes, error) {
	var taskLanes []Lane

	files, err := os.ReadDir(".projectSpy")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		for i, lane := range config.Lanes {
			if file.Name() == lane.Dir {
				taskLanes = append(taskLanes, Lane{
					Slug:  file.Name(),
					Name:  file.Name(),
					Title: config.Lanes[i].Name,
				})

				config.Lanes[i].HasDir = true
			}
		}
	}

	for i := 0; i < len(config.Lanes); i++ {
		j := slices.IndexFunc(taskLanes, func(lane Lane) bool {
			return lane.Slug == config.Lanes[i].Dir
		})

		if j == -1 {
			continue
		}

		taskLanes[j].Order = i
	}

	sort.SliceStable(taskLanes, func(i, j int) bool {
		return taskLanes[i].Order < taskLanes[j].Order
	})

	return taskLanes, nil
}

func ListTasks(lanes Lanes) {
	for k, lane := range lanes {
		entries, err := os.ReadDir(".projectSpy/" + lane.Slug)
		if err != nil {
			log.Fatal(err)
		}

		for _, e := range entries {
			task, err := parseFile(".projectSpy/" + lane.Slug + "/" + e.Name())
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

func SetupWatcher(eventBus *eventBus.EventBus[string], lanes Lanes) *fsnotify.Watcher {
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

				i := slices.IndexFunc(lanes, func(lane Lane) bool {
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

					task, err := parseFile(event.Name)
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

func RenderTaskLanes(config *config.Config, lanes []Lane) []Lane {
	for i := 0; i < len(lanes); i++ {
		lane := lanes[i]

		// Sort tasks by order ascending, then priority and modified time
		sort.SliceStable(lanes[i].Tasks, func(i, j int) bool {
			if lane.Tasks[i].Order != 0 && lane.Tasks[j].Order == 0 {
				// Items with Order > 0 come before items with Order == 0
				return true
			} else if lane.Tasks[i].Order == 0 && lane.Tasks[j].Order != 0 {
				// Items with Order == 0 come after items with Order > 0
				return false
			} else if lane.Tasks[i].Order != lane.Tasks[j].Order {
				// Within Order > 0, sort by original Order value
				return lane.Tasks[i].Order < lane.Tasks[j].Order
			} else if lane.Tasks[i].Order == 0 && lane.Tasks[j].Order == 0 {
				// Within Order == 0, sort by Priority and then Modified
				if lane.Tasks[i].Priority != lane.Tasks[j].Priority {
					return lane.Tasks[i].Priority > lane.Tasks[j].Priority
				}
				return lane.Tasks[j].ModifiedTime.Before(lane.Tasks[i].ModifiedTime)
			}

			return false
		})

		lanes[i].Tasks = lane.Tasks
		lanes[i].Count = len(lanes[i].Tasks)
	}

	return lanes
}
