package main

import (
	"net/http"

	"projectspy.dev/ui"
)

func (app *application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	mux.Handle("GET /", http.HandlerFunc(app.home))
	mux.Handle("GET /view/{lane}/{filename}", http.HandlerFunc(app.view))
	mux.Handle("GET /manifest.json", http.HandlerFunc(app.manifest))

	return mux
}
