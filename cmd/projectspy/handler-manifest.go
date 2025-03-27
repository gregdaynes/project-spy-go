package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) manifest(w http.ResponseWriter, r *http.Request) {
	type Icon struct {
		Src     string `json:"src"`
		Sizes   string `json:"sizes"`
		Type    string `json:"type"`
		Purpose string `json:"purpose"`
	}

	type Manifest struct {
		Name            string `json:"name"`
		ShortName       string `json:"short_name"`
		StartURL        string `json:"start_url"`
		Scope           string `json:"scope"`
		Icons           []Icon `json:"icons"`
		ThemeColor      string `json:"theme_color"`
		BackgroundColor string `json:"background_color"`
		Display         string `json:"display"`
		Orientation     string `json:"orientation"`
	}

	payload := Manifest{
		// TODO replace these with correct app name
		Name:      "ProjectSpy",
		ShortName: "ps",
		// TODO replace these with correct local url
		StartURL: "http://example.com",
		Scope:    "http://example.com",
		Icons: []Icon{
			// TODO do these even exist
			{Src: "/static/android-chrome-192x192.png", Sizes: "192x192", Type: "image/png"},
			{Src: "/static/android-chrome-512x512.png", Sizes: "512x512", Type: "image/png", Purpose: "any maskable"},
		},
		ThemeColor:      "#000000",
		BackgroundColor: "#000000",
		Display:         "standalone",
		Orientation:     "portait",
	}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonStr)
}
