package main

import (
	"httpserver/internal/config"
	"httpserver/internal/http-server/handlers/redirect"
	delete2 "httpserver/internal/http-server/handlers/url/delete"
	"httpserver/internal/http-server/handlers/url/save"
	"httpserver/internal/http-server/middleware/logger"
	"httpserver/internal/lib/logger/handlers/slogpretty"
	"httpserver/internal/lib/logger/sl"
	"httpserver/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
	cfgPath  = "config/local.yaml"
)

func main() {
	//initializing config and logger
	cfg := config.MustLoad(cfgPath)

	log := setupLogger(cfg.Env)

	log.Info("initializing server...", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	//initializing storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	storage_backup = storage.Clone()

	//initializing router and middleware
	router := chi.NewRouter()

	//middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//route handlers

	//adding alias
	router.Post("/url", save.New(log, storage))
	log.Info("starting server", slog.String("address", cfg.Address))

	//redirecting from existing alias
	router.Get("/{alias}", redirect.New(log, storage))

	//deleting alias
	router.Delete("/{alias}", delete2.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	//listening
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stopped")
}

// Setting up prettyLogger for local environment and json logger for dev and prod
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettyLogger()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// initializing prettyLogger
func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
