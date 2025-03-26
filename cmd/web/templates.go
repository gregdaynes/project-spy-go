package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"projectspy.dev/ui"
)

type templateData struct {
	message     string
	TaskLanes   map[int]ViewLaneModel
	CurrentTask ViewTaskModel
	ShowTask    bool
	ShowInfo    bool
	SearchData  string
}

type ViewLaneModel struct {
	Name     string
	Slug     string
	Tasks    map[string]ViewTaskModel
	Count    int
	Selected bool
}

type ViewTaskModel struct {
	Name            string
	ID              string
	Lane            string
	Title           string
	Body            string
	DescriptionHTML string
	Description     string
	Priority        int
	Tags            []string
	FullPath        string
	RelativePath    string
	Filename        string
	ModifiedTime    time.Time
	CreatedTime     time.Time
	Order           int
	ShowDetails     bool
	Actions         map[string]ViewActionModel
	AvailableLanes  map[string]ViewLaneModel
}

type ViewActionModel struct {
	Label  string
	Name   string
	Method string
	Action string
}

var functions = template.FuncMap{}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		message:   "Hello, world!",
		TaskLanes: make(map[int]ViewLaneModel),
	}
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		templateSet, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = templateSet
	}

	return cache, nil
}
