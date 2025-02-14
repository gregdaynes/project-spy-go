package main

import (
	"crypto/tls"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type application struct {
	debug  bool
	logger *slog.Logger
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

	app := &application{
		debug:  *debug,
		logger: logger,
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
	err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}
