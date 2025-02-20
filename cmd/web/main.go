package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
)

type application struct {
	debug         bool
	logger        *slog.Logger
	templateCache map[string]*template.Template
	taskLanes     map[string]TaskLane
	watcher       *fsnotify.Watcher
}

func main() {
	addr := flag.String("addr", ":8443", "HTTP network address")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	slogHandlerOptions := slog.HandlerOptions{}
	slogHandlerOptions.Level = slog.LevelInfo

	if *debug {
		slogHandlerOptions.AddSource = true
		slogHandlerOptions.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slogHandlerOptions))

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// setup task lanes
	taskLanes, err := newTaskLanes()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	watcher := setupWatcher(taskLanes)
	defer watcher.Close()

	listTasks(taskLanes)

	fmt.Println(taskLanes)

	app := &application{
		debug:         *debug,
		logger:        logger,
		templateCache: templateCache,
		taskLanes:     taskLanes,
		watcher:       watcher,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.Routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	logger.Info("starting server", slog.String("addr", *addr))
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
