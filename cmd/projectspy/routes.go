package main

import (
	"net/http"

	"projectspy.dev/web"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	if app.dev {
		fileServer := http.FileServer(http.Dir("./web/static/"))
		mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	} else {
		mux.Handle("GET /static/", http.FileServerFS(web.Files))
	}
	mux.Handle("GET /", http.HandlerFunc(app.home))
	mux.Handle("GET /file/{filename}", http.HandlerFunc(app.viewFile))
	mux.Handle("GET /{tid}", http.HandlerFunc(app.viewById))
	mux.Handle("GET /info", http.HandlerFunc(app.info))
	mux.Handle("POST /new", http.HandlerFunc(app.createTask))
	mux.Handle("GET /new", http.HandlerFunc(app.newTask))
	mux.Handle("GET /view/{lane}/{filename}", http.HandlerFunc(app.view))
	mux.Handle("POST /update/{lane}/{filename}", http.HandlerFunc(app.update))
	mux.Handle("GET /delete/{lane}/{filename}", http.HandlerFunc(app.delete))
	mux.Handle("POST /delete/{lane}/{filename}", http.HandlerFunc(app.deleteConfirm))
	mux.Handle("GET /archive/{lane}/{filename}", http.HandlerFunc(app.archive))
	mux.Handle("POST /archive/{lane}/{filename}", http.HandlerFunc(app.archiveConfirm))
	mux.Handle("GET /manifest.json", http.HandlerFunc(app.manifest))
	mux.Handle("POST /attach-file", http.HandlerFunc(app.attachFile))

	return mux
}
