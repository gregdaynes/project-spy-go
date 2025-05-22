package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/fsnotify/fsnotify"
	"projectspy.dev/internal/browser"
	"projectspy.dev/internal/config"
	eventBus "projectspy.dev/internal/event-bus"
	"projectspy.dev/internal/task"
	"projectspy.dev/web"
)

type application struct {
	debug         bool
	logger        *slog.Logger
	templateCache map[string]*template.Template
	taskLanes     []task.Lane
	watcher       *fsnotify.Watcher
	eventBus      *eventBus.EventBus[string]
	config        *config.Config
}

func main() {
	port := flag.String("port", "0", "HTTP network port to listen on")
	debug := flag.Bool("debug", false, "Enable debug mode")
	init := flag.Bool("init", false, "Initialize the project")
	flag.Parse()

	if *init {
		config.InitProject()
		os.Exit(0)
	}

	slogHandlerOptions := slog.HandlerOptions{}
	slogHandlerOptions.Level = slog.LevelInfo

	if *debug {
		slogHandlerOptions.AddSource = true
		slogHandlerOptions.Level = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slogHandlerOptions))

	cfg, err := config.NewConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	templateCache, err := web.NewTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// setup task lanes
	taskLanes, err := task.NewTaskLanes(&cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	eventBus := eventBus.NewEventBus[string]()

	watcher := task.SetupWatcher(eventBus, taskLanes)

	task.ListTasks(taskLanes)

	app := &application{
		debug:         *debug,
		logger:        logger,
		templateCache: templateCache,
		taskLanes:     taskLanes,
		watcher:       watcher,
		eventBus:      eventBus,
		config:        &cfg,
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	// logger.Info("starting server", slog.String("port", *port))
	l, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatal(err)
	}

	prt := fmt.Sprint(l.Addr().(*net.TCPAddr).Port)
	fmt.Printf("Project Spy is running\nhttp://localhost:%v", prt)

	err = browser.Open("http://localhost:" + prt)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Serve(l))
}

func (app *application) newTemplateData(_ *http.Request) web.TemplateData {
	return web.TemplateData{
		Message: "Hello, world!",
	}
}
