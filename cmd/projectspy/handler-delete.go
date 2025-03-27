package main

import (
	"log"
	"net/http"
	"os"
)

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	task, ok := app.taskLanes[lane].Tasks[filename]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// get the file path
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := cwd + "/.projectSpy" + task.RelativePath

	ch := make(chan int)

	waitForWrite := func(event string) {
		if event == lane+"/"+filename {
			ch <- 1
		}
	}

	go func() {
		app.eventBus.Subscribe("remove", &waitForWrite)
		err := os.Remove(path)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("remove", &waitForWrite)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
