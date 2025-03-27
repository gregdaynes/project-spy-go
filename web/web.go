package web

import (
	"embed"
	"text/template"

	"projectspy.dev/internal/task"
)

//go:embed "html" "static"
var Files embed.FS

type TemplateData struct {
	Message     string
	TaskLanes   map[int]task.TaskLane
	CurrentTask task.Task
	ShowTask    bool
	ShowInfo    bool
	SearchData  string
}

var functions = template.FuncMap{}

func NewTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	patterns := []string{
		"html/base.tmpl",
		"html/partials/*.tmpl",
	}

	templateSet, err := template.New("app").Funcs(functions).ParseFS(Files, patterns...)
	if err != nil {
		return nil, err
	}

	cache["app"] = templateSet

	return cache, nil
}
