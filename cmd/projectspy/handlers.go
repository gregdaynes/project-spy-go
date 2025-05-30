package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"projectspy.dev/internal/search"
	"projectspy.dev/internal/task"
	"projectspy.dev/internal/util"
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
	data.SearchData = search.Data(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)
	data.ShowConfirm = true
	data.Confirm = web.Confirm{
		Title: "Archive",
		Body:  "Are you sure you want to archive task <samp>" + t.Title + "</samp>?",
	}
	data.Confirm.Actions = make(map[string]task.Action)
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

	slug := strings.TrimSuffix(t.Filename, ".md")
	filename := slug

	// check if new file already exists, if so, append an incrementing number to the filename
	i := 1
	for {
		_, err := os.Stat(filepath("_archive", filename+".md"))
		if err != nil {
			break
		}
		filename = slug + "-" + strconv.Itoa(i)
		i++
	}
	filename += ".md"

	ch := make(chan int)
	waitForWrite := wait(ch, t.RelativePath)

	go func() {
		app.eventBus.Subscribe("remove", "archiveRemove", waitForWrite)
		err := os.Rename(filepath(t.Lane, t.Filename), filepath("_archive", filename))
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("remove", "archiveRemove")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) createTask(w http.ResponseWriter, r *http.Request) {
	qName := r.FormValue("name")
	qContent := r.FormValue("content")
	qLane := r.FormValue("lane")

	slug := slug.Make(qName)
	uId := util.GenerateId(slug)

	content := qName + "\n===\n\n" + qContent
	content = appendChangelog(content, "Created task")

	filename := slug + ":" + uId

	// check if file already exists, if so, append an incrementing number to the filename
	i := 1
	for {
		_, err := os.Stat(filepath(qLane, filename+".md"))
		if err != nil {
			break
		}
		filename = slug + "-" + strconv.Itoa(i)
		i++
	}

	filename += ".md"

	path := filepath(qLane, filename)

	ch := make(chan int)
	waitForWrite := wait(ch, qLane+"/"+filename)

	go func() {
		app.eventBus.Subscribe("update", "createTask", waitForWrite)
		err := os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("update", "createTask")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	t, ok := app.getTask(qLane, qFile)
	if !ok {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := app.newTemplateData(r)
	data.SearchData = search.Data(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)
	data.ShowConfirm = true
	data.Confirm = web.Confirm{
		Title: "Delete",
		Body:  "Are you sure you want to delete task <samp>" + t.Title + "</samp>?",
	}
	data.Confirm.Actions = make(map[string]task.Action)
	data.Confirm.Actions["Confirm"] = task.Action{
		Label:  "Delete",
		Name:   "delete",
		Method: "POST",
		Action: "/delete/" + r.PathValue("lane") + "/" + r.PathValue("filename"),
	}

	data.Confirm.Actions["Close"] = task.Action{
		Label:  "Close",
		Name:   "close",
		Method: "GET",
		Action: "/view/" + t.RelativePath,
	}

	app.render(w, r, http.StatusOK, data)
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
		app.eventBus.Subscribe("remove", "deleteTask", waitForWrite)
		err := os.Remove(path)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch
	app.eventBus.Unsubscribe("remove", "deleteTask")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.SearchData = search.Data(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	app.render(w, r, http.StatusOK, data)
}

func (app *application) info(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.SearchData = search.Data(app.taskLanes)
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
	var qLane string

	if _, ok := r.URL.Query()["lane"]; ok {
		qLane = r.URL.Query()["lane"][0]
	}

	data := app.newTemplateData(r)
	data.SearchData = search.Data(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	data.CurrentTask = task.Task{
		Title:    "",
		Priority: 0,
		Tags:     []string{},
		Lane:     qLane,
		Lanes:    app.taskLanes,
	}
	data.ShowTask = true

	app.render(w, r, http.StatusOK, data)
}

func (app *application) update(w http.ResponseWriter, r *http.Request) {
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	oldPath := filepath(qLane, qFile)
	_, err := os.ReadFile(oldPath)
	if err != nil {
		log.Fatal("file not found")
	}

	content := r.FormValue("content")
	content = scrub(content)
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

	filename := strings.TrimSuffix(qFile, ".md")
	var path string

	if newLane != qLane {
		// check file already exists, if so, append an incrementing number to the filename
		i := 1
		for {
			_, err := os.Stat(filepath(newLane, filename+".md"))
			if err != nil {
				break
			}
			filename = filename + "-" + strconv.Itoa(i)
			i++
		}

		err := os.Remove(oldPath)
		if err != nil {
			log.Fatal(err)
		}
		content = appendChangelog(content, "Moved task from "+qLane+" to "+newLane)
		qLane = newLane
	}

	filename += ".md"
	path = filepath(newLane, filename)

	ch := make(chan int)
	waitForWrite := wait(ch, newLane+"/"+filename)

	go func() {
		app.eventBus.Subscribe("update", "update", waitForWrite)
		err = os.WriteFile(path, []byte(content), 0644)
	}()

	<-ch
	app.eventBus.Unsubscribe("update", "update")

	t, ok = app.getTask(newLane, filename)
	if !ok {
		log.Fatal("task not found")
	}

	http.Redirect(w, r, "/view/"+t.RelativePath, http.StatusSeeOther)
}

func (app *application) view(w http.ResponseWriter, r *http.Request) {
	qLane := r.PathValue("lane")
	qFile := r.PathValue("filename")

	data := app.newTemplateData(r)
	data.SearchData = search.Data(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	t, ok := app.getTask(qLane, qFile)
	if !ok {
		data.ShowError = true
		data.ErrorDialog = web.ErrorDialog{
			Title: "Task not found",
			Body:  "The task you are looking for could not be found.",
		}
		app.render(w, r, http.StatusNotFound, data)
		return
	}

	data.CurrentTask = t
	data.CurrentTask.Lanes = app.taskLanes

	data.ShowTask = true

	app.render(w, r, http.StatusOK, data)
}

func (app *application) viewById(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.SearchData = search.Data(app.taskLanes)
	data.TaskLanes = task.RenderTaskLanes(app.config, app.taskLanes)

	tid := r.PathValue("tid")
	t, ok := app.getTaskById(tid)
	if !ok {
		data.ShowError = true
		data.ErrorDialog = web.ErrorDialog{
			Title: "Task not found",
			Body:  "The task you are looking for could not be found.",
		}
		app.render(w, r, http.StatusNotFound, data)
		return
	}

	data.CurrentTask = t
	data.CurrentTask.Lanes = app.taskLanes

	data.ShowTask = true

	app.render(w, r, http.StatusOK, data)
}

func (app *application) getTask(lane, filename string) (t task.Task, ok bool) {
	i := slices.IndexFunc(app.taskLanes, func(l task.Lane) bool {
		return l.Name == lane
	})
	if i == -1 {
		return task.Task{}, false
	}

	j := slices.IndexFunc(app.taskLanes[i].Tasks, func(task task.Task) bool {
		return task.Filename == filename
	})
	if j == -1 {
		return task.Task{}, false
	}

	return app.taskLanes[i].Tasks[j], true
}

func (app *application) getTaskById(tid string) (t task.Task, ok bool) {
	var iL, iT int

	for i, lane := range app.taskLanes {
		iT = slices.IndexFunc(lane.Tasks, func(t task.Task) bool {
			return t.ID == tid
		})

		if iT != -1 {
			iL = i
			break
		}
	}

	if iT < 0 {
		return task.Task{}, false
	}

	return app.taskLanes[iL].Tasks[iT], true
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
