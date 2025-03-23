package main

import (
	"crypto/tls"
	"flag"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
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

	listTasks(taskLanes)

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
	l, err := net.Listen("tcp", ":8443")
	if err != nil {
		log.Fatal(err)
	}

	err = open("https://localhost:8443")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.ServeTLS(l, "./tls/cert.pem", "./tls/key.pem"))
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
