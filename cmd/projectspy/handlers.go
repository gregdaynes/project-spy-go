package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"projectspy.dev/internal/search"
	"projectspy.dev/internal/task"
	"projectspy.dev/web"
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
func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	task2, ok := app.taskLanes[lane].Tasks[filename]
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// get the file path
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := cwd + "/.projectSpy" + task2.RelativePath

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
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
func (app *application) info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	data.ShowInfo = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
func (app *application) manifest(w http.ResponseWriter, r *http.Request) {
	type icon struct {
		Src     string `json:"src"`
		Sizes   string `json:"sizes"`
		Type    string `json:"type"`
		Purpose string `json:"purpose"`
	}

	type manifest struct {
		Name            string `json:"name"`
		ShortName       string `json:"short_name"`
		StartURL        string `json:"start_url"`
		Scope           string `json:"scope"`
		Icons           []icon `json:"icons"`
		ThemeColor      string `json:"theme_color"`
		BackgroundColor string `json:"background_color"`
		Display         string `json:"display"`
		Orientation     string `json:"orientation"`
	}

	payload := manifest{
		// TODO replace these with correct app name
		Name:      "ProjectSpy",
		ShortName: "ps", // TODO replace these with correct local url
		StartURL:  "https://example.com",
		Scope:     "https://example.com",
		Icons: []icon{
			// TODO do these even exist
			{Src: "/static/android-chrome-192x192.png", Sizes: "192x192", Type: "image/png"},
			{Src: "/static/android-chrome-512x512.png", Sizes: "512x512", Type: "image/png", Purpose: "any maskable"},
		},
		ThemeColor:      "#000000",
		BackgroundColor: "#000000",
		Display:         "standalone",
		Orientation:     "portrait",
	}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonStr)
	if err != nil {
		app.serverError(w, r, err)
	}
}
func (app *application) newTask(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	newTask := task.Task{}

	data.CurrentTask = web.ViewTaskModel{
		Title:          "",
		Body:           "",
		ShowDetails:    true,
		Priority:       0,
		Tags:           []string{},
		AvailableLanes: task.GetAvailableLanes(&newTask, app.taskLanes),
		Actions:        task.GetAvailableActions(&newTask),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
func (app *application) update(w http.ResponseWriter, r *http.Request) {
	// 1. receive updated content
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	path := cwd + "/.projectSpy/" + lane + "/" + filename
	fmt.Println("path:", path)
	_, err = os.ReadFile(path)
	if err != nil {
		log.Fatal("file not found")
	}

	content := r.FormValue("content")
	newLane := r.FormValue("lane")

	contentHasher := sha1.New()
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.TrimSpace(content)

	contentHasher.Write([]byte(content))
	contentHash := hex.EncodeToString(contentHasher.Sum(nil))

	rawContents := app.taskLanes[lane].Tasks[filename].RawContents
	rawContents = strings.ReplaceAll(rawContents, "\r\n", "\n")
	rawContents = strings.TrimSpace(rawContents)

	rawHasher := sha1.New()
	rawHasher.Write([]byte(rawContents))
	rawHash := hex.EncodeToString(rawHasher.Sum(nil))

	if contentHash == rawHash && lane == newLane {
		http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
		return
	}

	if contentHash != rawHash {
		// get last line of content
		last := strings.Split(content, "\n")[len(strings.Split(content, "\n"))-1]
		// regexto detect if it is a changelog entry
		re := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}\t.*`)
		if re.MatchString(last) {
			// create "yyyy-mm-dd entry" timestamp
			timestamp := time.Now().Format("2006-01-02 15:04")
			content += "\n" + timestamp + "\tUpdated task"
		} else {
			timestamp := time.Now().Format("2006-01-02 15:04")
			content += "\n\n---\n\n" + timestamp + "\tUpdated task"
		}
	}

	if newLane != lane {
		oldLane := lane
		oldPath := path
		path = cwd + "/.projectSpy/" + newLane + "/" + filename
		// send message to delete old file
		err := os.Remove(oldPath)
		if err != nil {
			log.Fatal(err)
		}
		lane = newLane
		last := strings.Split(content, "\n")[len(strings.Split(content, "\n"))-1]
		re := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}\t.*`)
		if re.MatchString(last) {
			timestamp := time.Now().Format("2006-01-02 15:04")
			content += "\n" + timestamp + "\tMoved task from " + oldLane + " to " + newLane
		} else {
			timestamp := time.Now().Format("2006-01-02 15:04")
			content += "\n\n---\n\n" + timestamp + "\tMoved task from " + oldLane + " to " + newLane
		}
	}

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
func (app *application) view(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	data := app.newTemplateData(r)

	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	currentTask := app.taskLanes[lane].Tasks[filename]

	data.CurrentTask = web.ViewTaskModel{
		Title:          currentTask.Title,
		Body:           currentTask.RawContents,
		ShowDetails:    currentTask.HasPriorityOrTags(),
		Priority:       currentTask.Priority,
		Tags:           currentTask.Tags,
		AvailableLanes: task.GetAvailableLanes(&currentTask, app.taskLanes),
		Actions:        task.GetAvailableActions(&currentTask),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, "home.tmpl", data)
}
