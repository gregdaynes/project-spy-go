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

	if !app.taskExists(lane, filename) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	currentPath := filepath(lane, filename)
	newPath := filepath("_archive", filename)

	ch := make(chan int)
	waitForWrite := app.createWaitForWriteFunc(ch, lane, filename)

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
	name := r.FormValue("name")
	content := r.FormValue("content")
	lane := r.FormValue("lane")

	filename := slug.Make(name) + ".md"
	content = name + "\n===\n\n" + content
	path := app.getTaskPath(lane, filename)

	ch := make(chan int)
	waitForWrite := app.createWaitForWriteFunc(ch, lane, filename)

	go func() {
		app.eventBus.Subscribe("update", &waitForWrite)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("update", &waitForWrite)

	http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	if !app.taskExists(lane, filename) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	path := app.getTaskPath(lane, filename)

	ch := make(chan int)
	waitForWrite := app.createWaitForWriteFunc(ch, lane, filename)

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
		Name:      "ProjectSpy",
		ShortName: "ps",
		StartURL:  "https://example.com",
		Scope:     "https://example.com",
		Icons: []icon{
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
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	path := app.getTaskPath(lane, filename)
	_, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("file not found")
	}

	content := r.FormValue("content")
	newLane := r.FormValue("lane")

	currentContent := app.taskLanes[lane].Tasks[filename].RawContents

	if same(content, currentContent) && lane == newLane {
		http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
		return
	}

	content = appendChangelog(content, "Updated task")

	if newLane != lane {
		oldPath := path
		path = app.getTaskPath(newLane, filename)
		err := os.Remove(oldPath)
		if err != nil {
			log.Fatal(err)
		}
		content = appendChangelog(content, "Moved task from "+lane+" to "+newLane)
		lane = newLane
	}

	ch := make(chan int)
	waitForWrite := app.createWaitForWriteFunc(ch, lane, filename)

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

func (app *application) taskExists(lane, filename string) bool {
	_, ok := app.taskLanes[lane].Tasks[filename]

	return ok
}

func filepath(lane, filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cwd + "/.projectSpy/" + lane + "/" + filename
}

func (app *application) getTaskPath(lane, filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cwd + "/.projectSpy/" + lane + "/" + filename
}

func (app *application) createWaitForWriteFunc(ch chan int, lane, filename string) func(string) {
	return func(event string) {
		if event == lane+"/"+filename {
			ch <- 1
		}
	}
}

func same(a, b string) bool {
	fmt.Println(hash(scrub(a)), hash(scrub(b)))
	return hash(scrub(a)) == hash(scrub(b))
}

func scrub(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimSpace(s)
	return s
}

func hash(s string) string {
	H := sha1.New()
	H.Write([]byte(s))
	h := hex.EncodeToString(H.Sum(nil))
	return h
}

func appendChangelog(content, change string) string {
	last := strings.Split(content, "\n")[len(strings.Split(content, "\n"))-1]
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}\t.*`)
	timestamp := time.Now().Format("2006-01-02 15:04")

	if re.MatchString(last) {
		content += "\n" + timestamp + "\t" + change
	} else {
		content += "\n\n---\n\n" + timestamp + "\t" + change
	}

	return content
}
