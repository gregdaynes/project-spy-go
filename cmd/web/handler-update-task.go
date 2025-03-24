package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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
	fmt.Println(path)
	_, err = os.ReadFile(path)
	if err != nil {
		log.Fatal("file not found")
	}

	content := r.FormValue("content")
	err = os.WriteFile(path, []byte(content), 0644)

	http.Redirect(w, r, "/view/"+lane+"/"+filename, http.StatusSeeOther)
}
