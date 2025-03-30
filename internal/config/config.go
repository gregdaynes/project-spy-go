package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigLane struct {
	Dir    string `json:"dir"`
	Name   string `json:"name"`
	HasDir bool
}
type Config struct {
	Lanes []ConfigLane `json:"lanes"`
}

func NewConfiguration() (config Config, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.ReadFile(cwd + "/.projectSpy/projectspy.json")
	if err != nil {
		config = Config{
			Lanes: []ConfigLane{
				{
					Dir:  "inbox",
					Name: "Inbox",
				},
				{
					Dir:  "backlog",
					Name: "Backlog",
				},
				{
					Dir:  "blocked",
					Name: "Blocked",
				},
				{
					Dir:  "in-progress",
					Name: "In Progress",
				},
				{
					Dir:  "done",
					Name: "Done",
				},
			},
		}
		return config, nil
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}
