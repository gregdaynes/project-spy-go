package main

import "net/http"

func (app *application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /ping/", http.HandlerFunc(ping))

	return mux
}
