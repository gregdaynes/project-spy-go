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
	TaskLanes   []task.Lane
	CurrentTask task.Task
	ShowTask    bool
	ShowInfo    bool
	SearchData  string
	ShowConfirm bool
	Confirm     Confirm
	ShowError   bool
	ErrorDialog ErrorDialog
}

type Confirm struct {
	Title   string
	Body    string
	Actions map[string]task.Action
}

type ErrorDialog struct {
	Title string
	Body  string
}

var functions = template.FuncMap{}

func NewTemplateCache(devMode *bool) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	var templateSet *template.Template
	var err error

	if *devMode {
		files := []string{
			"./web/html/base.tmpl",
			"./web/html/partials/confirm-dialog.tmpl",
			"./web/html/partials/error-dialog.tmpl",
			"./web/html/partials/filter.tmpl",
			"./web/html/partials/info-dialog.tmpl",
			"./web/html/partials/lane.tmpl",
			"./web/html/partials/lanes.tmpl",
			"./web/html/partials/new-task.tmpl",
			"./web/html/partials/panel-about.tmpl",
			"./web/html/partials/panel-info.tmpl",
			"./web/html/partials/task-card.tmpl",
			"./web/html/partials/task.tmpl",
		}
		templateSet, err = template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
	} else {
		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
		}

		templateSet, err = template.New("app").Funcs(functions).ParseFS(Files, patterns...)
		if err != nil {
			return nil, err
		}
	}

	cache["app"] = templateSet

	return cache, nil
}
