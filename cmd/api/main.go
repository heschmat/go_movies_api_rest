package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

// Application version
//@TODO: later generate this automatically at build time.
const version = "1.0.0"

// The *config* struct holds all the configuration settings for the application.
// We will read in *configuration settings* from the command-line flags when the application starts.
// port: the network port that we want the server to listen on
// env : the operating environment for the application (development, staging, production)
// ...
type config struct {
	port int
	env  string
}

// The *application* struct holds all the `dependencies` for the HTTP handlers, helpers & middleware.
type application struct {
	config config
	logger *slog.Logger
}


func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Inisitalize a new structured logger.
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,  		// the filename & line number of the calling source code
		Level: slog.LevelDebug,	// the minimum log level
	}))

	// Declare an instance of the application struct.
	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// Start the HTTP server.
	logger.Info("Starting server", "addr", srv.Addr, "env", cfg.env)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
