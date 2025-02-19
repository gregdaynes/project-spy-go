package main

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"

	"projectspy.dev/ui"
)

type templateData struct {
	message   string
	TaskLanes TaskLanes
}

var functions = template.FuncMap{}

func (app *application) newTemplateData(r *http.Request) templateData {
	fmt.Println(r)

	return templateData{
		message:   "Hello, world!",
		TaskLanes: app.taskLanes,
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
			// "html/partials/*.tmpl",
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
