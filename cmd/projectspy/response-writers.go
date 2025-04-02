package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, data any) {
	templateSet, ok := app.templateCache["app"]

	if !ok {
		err := fmt.Errorf("the template %s does not exist", "app")
		app.serverError(w, r, err)
		return
	}

	buf := new(bytes.Buffer)

	err := templateSet.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	if err != nil {
		return
	}
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri), slog.Any("trace", trace))

	if app.debug {
		body := fmt.Sprintf("%s\n%s", err, trace)
		http.Error(w, body, http.StatusInternalServerError)
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
