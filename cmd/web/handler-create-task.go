package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gosimple/slug"
)

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	// create task receives a form input
	// parses / validates the form input
	// creates a new file in the path
	name := r.FormValue("name")
	content := r.FormValue("content")
	lane := r.FormValue("lane")

	// slugify name and add extension
	filename := slug.Make(name) + ".md"

	// create the task by
	content = name + "\n===\n\n" + content

	// create the path
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := cwd + "/.projectSpy/" + lane + "/" + filename

	ch := make(chan int)

	waitForWrite := func(event string) {
		if event == lane+"/"+filename {
			ch <- 1
		}
	}

	go func() {
		app.eventBus.Subscribe("update", &waitForWrite)
		err = os.WriteFile(path, []byte(content), 0644)
	}()

	<-ch
	app.eventBus.Unsubscribe("update", &waitForWrite)

	http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
}
