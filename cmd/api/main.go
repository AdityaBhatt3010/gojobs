package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/karnop/gojobs/internal/data"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// defining a struct to hold application dependencies
// it makes the handlers cleaner because they can access the DB via this struct
type application struct {
	DB     *sql.DB
	Users  data.UserModel
	Applications data.JobApplicationModel
	Jobs data.JobModel
	Logger *slog.Logger
}

// entry point of the application
func main() {
	// initializing structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// loading the env file
	err := godotenv.Load()
	if err != nil {
		logger.Info("No .env file found, relying on system environment variables")
    	os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
    	logger.Error("PORT environment variable not set")
    	os.Exit(1)
	}

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
    	logger.Error("DB_DSN environment variable not set")
    	os.Exit(1)
	}

logger.Info("Connecting to Database...")

// calling helper function to open the connection
db, err := openDB(dsn)
if err != nil {
    logger.Error("Cannot connect to database", "error", err)
    os.Exit(1)
}

defer db.Close()
logger.Info("SUCCESS! Database connection established")

	app := &application{
		DB:     db,
		Users:  data.UserModel{DB: db},
		Applications: data.JobApplicationModel{DB: db},
		Jobs : data.JobModel{DB: db},
		Logger: logger,
	}

	// NewServeMux is a request multiplier (router).
	// It matches the URL of incoming request against a list of registered patterns
	// and calls the corresponding handler.
	mux := http.NewServeMux()

	// Register a simple health check route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the GoJobs API")
	})

	mux.HandleFunc("GET /jobs", app.listJobsHandler)
	mux.HandleFunc("POST /jobs", app.authenticate(app.createJobHandler))
	mux.HandleFunc("GET /jobs/{id}", app.getJobHandler)
	mux.HandleFunc("POST /users", app.registerUserHandler)
	mux.HandleFunc("POST /users/login", app.loginUserHandler)
	mux.HandleFunc("POST /jobs/{id}/apply", app.authenticate(app.applyJobHandler))

	logger.Info("Starting server", "addr", port, "env", "development")

	// defining the server struct
	srv := &http.Server{
		Addr:  ":8080",
		Handler: app.enableCORS(mux), 
		IdleTimeout: time.Minute,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 30*time.Second,
	}

	// creating a specific channel to listen for shutdown signals
	// buffered channel of size 1
	shutdownError := make(chan error)

	// starting the server in a background routine
	go func() {
		logger.Info("Starting server", "addr", srv.Addr, "env", "development")
		err := srv.ListenAndServe()

		// ListenAndServe always returns a non nil error
		// If its just - ServerClosed (which happens when we shutdown), that's normal
		if !errors.Is(err, http.ErrServerClosed) {
			shutdownError <- err // sending real errors to the main thread
		}
	}()

	// listening for os signals
	// we want to catch interrupt(ctrl + c) and SIGTERM(docker/kubernetes stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// blocking the main thread
	// The code stops here and waits until we receive a signal or a server error
	select {
	case err := <-shutdownError:
		logger.Error("Server error", "error", err)
		os.Exit(1)

	case sig:= <-quit:
		logger.Info("Shutting down server", "signal", sig.String())
	}

	// graceful shutdown
	// we give active requests 5 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown stops accepting new requests and waits for active ones
	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
        err = srv.Close() // force close
	}

	logger.Info("Server stopped")  
}
