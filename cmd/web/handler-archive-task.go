package main

import (
	"log"
	"net/http"
	"os"
)

func (app *application) archive(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	_, ok := app.taskLanes[lane].Tasks[filename]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// get the file path
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	currentPath := cwd + "/.projectSpy/" + lane + "/" + filename
	newPath := cwd + "/.projectSpy/_archive/" + filename

	ch := make(chan int)

	waitForWrite := func(event string) {
		if event == lane+"/"+filename {
			ch <- 1
		}
	}

	go func() {
		app.eventBus.Subscribe("remove", &waitForWrite)
		err := os.Rename(currentPath, newPath)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("remove", &waitForWrite)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
