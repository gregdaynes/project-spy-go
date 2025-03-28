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
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	t, ok := app.getTask(qLane, qFile)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)
	data.ShowConfirm = true
	data.Confirm = web.Confirm{
		Title: "Archive",
		Body:  "Are you sure you want to archive task <samp>" + t.Title + "</samp>?",
	}
	data.Confirm.Actions = make(map[string]task.Action, 0)
	data.Confirm.Actions["Confirm"] = task.Action{
		Label:  "Archive",
		Name:   "archive",
		Method: "POST",
		Action: "/archive/" + r.PathValue("lane") + "/" + r.PathValue("filename"),
	}

	data.Confirm.Actions["Close"] = task.Action{
		Label:  "Close",
		Name:   "close",
		Method: "GET",
		Action: "/view/" + t.RelativePath,
	}

	app.render(w, r, http.StatusOK, data)
}

func (app *application) archiveConfirm(w http.ResponseWriter, r *http.Request) {
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	t, ok := app.getTask(qLane, qFile)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println("will wait for", t.RelativePath)
	ch := make(chan int)
	waitForWrite := wait(ch, t.RelativePath)

	go func() {
		app.eventBus.Subscribe("remove", &waitForWrite)
		err := os.Rename(filepath(t.Lane, t.Filename), filepath("_archive", t.Filename))
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("remove", &waitForWrite)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	qName := r.FormValue("name")
	qContent := r.FormValue("content")
	qLane := r.FormValue("lane")

	filename := slug.Make(qName) + ".md"
	content := qName + "\n===\n\n" + qContent
	path := filepath(qLane, filename)

	content = appendChangelog(content, "Created task")

	ch := make(chan int)
	waitForWrite := wait(ch, qLane+"/"+filename)

	go func() {
		app.eventBus.Subscribe("update", &waitForWrite)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("update", &waitForWrite)

	t, ok := app.getTask(qLane, filename)
	if !ok {
		log.Fatal("task not found")
	}

	http.Redirect(w, r, "/view/"+t.RelativePath, http.StatusSeeOther)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteConfirm(w http.ResponseWriter, r *http.Request) {
	lane := r.PathValue("lane")
	filename := r.PathValue("filename")

	t, ok := app.getTask(lane, filename)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	path := filepath(lane, filename)

	ch := make(chan int)
	waitForWrite := wait(ch, t.RelativePath)

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

	app.render(w, r, http.StatusOK, data)
}

func (app *application) info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)
	data.ShowInfo = true
	app.render(w, r, http.StatusOK, data)
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
	data.CurrentTask = task.Task{
		Title:          "",
		Body:           "",
		ShowDetails:    true,
		Priority:       0,
		Tags:           []string{},
		AvailableLanes: task.GetAvailableLanes(&newTask, app.taskLanes),
		Actions:        task.GetAvailableActions(&newTask, "create"),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, data)
}

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	path := filepath(qLane, qFile)
	_, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("file not found")
	}

	content := r.FormValue("content")
	newLane := r.FormValue("lane")
	t, ok := app.getTask(qLane, qFile)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if same(content, t.RawContents) && qLane == newLane {
		http.Redirect(w, r, "/view/"+t.RelativePath, http.StatusSeeOther)
		return
	}

	content = appendChangelog(content, "Updated task")

	if newLane != qLane {
		oldPath := path
		path = filepath(newLane, qFile)
		err := os.Remove(oldPath)
		if err != nil {
			log.Fatal(err)
		}
		content = appendChangelog(content, "Moved task from "+qLane+" to "+newLane)
		qLane = newLane
	}

	ch := make(chan int)
	waitForWrite := wait(ch, qLane+"/"+qFile)

	go func() {
		app.eventBus.Subscribe("update", &waitForWrite)
		err = os.WriteFile(path, []byte(content), 0644)
	}()

	<-ch
	app.eventBus.Unsubscribe("update", &waitForWrite)

	t, ok = app.getTask(qLane, qFile)
	if !ok {
		log.Fatal("task not found")
	}

	http.Redirect(w, r, "/view/"+t.RelativePath, http.StatusSeeOther)
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	t, ok := app.getTask(qLane, qFile)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.SearchData = search.SearchData(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	data.CurrentTask = task.Task{
		Title:          t.Title,
		Body:           t.RawContents,
		ShowDetails:    t.HasPriorityOrTags(),
		Priority:       t.Priority,
		Tags:           t.Tags,
		AvailableLanes: task.GetAvailableLanes(&t, app.taskLanes),
		Actions:        task.GetAvailableActions(&t, "edit"),
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, data)
}

func (app *application) getTask(lane, filename string) (t task.Task, ok bool) {
	t, ok = app.taskLanes[lane].Tasks[filename]

	return t, ok
}

func filepath(lane, filename string) string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	return cwd + "/.projectSpy/" + lane + "/" + filename
}

func wait(ch chan int, path string) func(string) {
	return func(event string) {
		if event == path {
			ch <- 1
		}
	}
}

func same(a, b string) bool {
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

	content = strings.TrimSpace(content)

	if re.MatchString(last) {
		content += "\n" + timestamp + "\t" + change
	} else {
		content += "\n\n---\n\n" + timestamp + "\t" + change
	}

	return content
}
