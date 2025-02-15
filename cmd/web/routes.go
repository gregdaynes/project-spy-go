package main

import (
	"net/http"

	"projectspy.dev/ui"
)

func (app *application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.Handle("GET /", http.HandlerFunc(app.home))
	mux.Handle("GET /ping/", http.HandlerFunc(ping))

	return mux
}
