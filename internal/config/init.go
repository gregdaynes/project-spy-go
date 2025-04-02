package config

import (
	"fmt"
	"os"
)

func InitProject() {
	os.MkdirAll(".projectSpy/backlog", os.ModePerm)
	os.MkdirAll(".projectSpy/in-progress", os.ModePerm)
	os.MkdirAll(".projectSpy/done", os.ModePerm)

	fmt.Println("Project initialized")
}
