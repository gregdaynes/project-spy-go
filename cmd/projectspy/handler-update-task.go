package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

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

	rawcontents := app.taskLanes[lane].Tasks[filename].RawContents
	rawcontents = strings.ReplaceAll(rawcontents, "\r\n", "\n")
	rawcontents = strings.TrimSpace(rawcontents)

	rawHasher := sha1.New()
	rawHasher.Write([]byte(rawcontents))
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
		os.Remove(oldPath)
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
