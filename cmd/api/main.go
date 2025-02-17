package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	// Import the pq driver so that it can register itself with *database/sql* package.
	"github.com/heschmat/go_movies_api_rest/internal/data"
	_ "github.com/lib/pq"
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
	db	 struct {
		dsn 			string			// connection string
		// The following fields hold the configuration settings for the connection pool.
		maxOpenConns	int
		maxIdleConns	int
		maxIdleTime		time.Duration 	//300ms, 4s, 5h27m
	}
}

// The *application* struct holds all the `dependencies` for the HTTP handlers, helpers & middleware.
type application struct {
	config config
	logger *slog.Logger
	models data.Models
}


func main() {
	// Declare an instance of the config struct.
	var cfg config

	// Read the command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	// Default to using the development DSN if no flag is provided.
	// sample dsn: "postgres://<user>:<password>@localhost/<db>"
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MOVIESDB_DSN"), "PostgreSQL DSN")

	// Read the connection pool settings.
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15 * time.Minute, "PostgreSQL max connection idle time")
	flag.Parse()

	// Inisitalize a new structured logger --------------------- //
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,  		// the filename & line number of the calling source code
		Level: slog.LevelDebug,	// the minimum log level
	}))

	// Create the connection pool ------------------------------ //
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	// Make sure the connection ppol is closed before the main() function exits.
	defer db.Close()

	logger.Info("database connection pool established")

	// Declare an instance of the application struct.
	app := &application{
		config: cfg,
		logger: logger,
		// Initialize a Models struct; passing in the connection pool as a parameter.
		models: data.NewModels(db),
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
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}


func openDB(cfg config) (*sql.DB, error) {
	// Create an empty connection pool.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	// Create a context with a 5-sec timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// If the connection couldn't be established successfully within the 5 second deadline,
	// then pint fails & return an error.
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
