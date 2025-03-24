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

	contentHasher := sha1.New()
	content := r.FormValue("content")
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

	if contentHash == rawHash {
		http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
		return
	}

	// get last line of content
	last := strings.Split(content, "\n")[len(strings.Split(content, "\n"))-1]
	// regexto detect if it is a changelog entry
	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}\s.*`)
	if re.MatchString(last) {
		// create "yyyy-mm-dd entry" timestamp
		timestamp := time.Now().Format("2006-01-02 15:04")
		content += "\n" + timestamp + "\tUpdated task"
	} else {
		timestamp := time.Now().Format("2006-01-02 15:04")
		content += "\n\n---\n\n" + timestamp + "\tUpdated task"
	}

	ch := make(chan int)

	waitForWrite := func(event string) {
		ch <- 1
	}

	go func() {
		app.eventBus.Subscribe("update", &waitForWrite)
		err = os.WriteFile(path, []byte(content), 0644)
	}()

	<-ch
	app.eventBus.Unsubscribe("update", &waitForWrite)

	http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
}
