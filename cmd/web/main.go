package main

import (
	"crypto/tls"
	"encoding/json"
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
	eventBus      *EventBus[string]
	config        *Config
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

	config, err := newConfiguration()
	if err != nil {
		log.Fatal(err)
	}

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

	eventBus := NewEventBus[string]()

	watcher := setupWatcher(eventBus, taskLanes)

	listTasks(taskLanes)

	app := &application{
		debug:         *debug,
		logger:        logger,
		templateCache: templateCache,
		taskLanes:     taskLanes,
		watcher:       watcher,
		eventBus:      eventBus,
		config:        &config,
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

	//err = open("https://localhost:8443")
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Fatal(srv.Serve(l))
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

type ConfigLane struct {
	Dir  string `json:"dir"`
	Name string `json:"name"`
}
type Config struct {
	Lanes []ConfigLane `json:"lanes"`
}

func newConfiguration() (config Config, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.ReadFile(cwd + "/.projectSpy/projectspy.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config, nil
}
