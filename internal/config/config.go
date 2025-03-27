package config

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigLane struct {
	Dir  string `json:"dir"`
	Name string `json:"name"`
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
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}
