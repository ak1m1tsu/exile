package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/effective-mobile-test-task/internal/config"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/sl"
	"github.com/romankravchuk/effective-mobile-test-task/internal/log"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/handlers/person/delete"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/handlers/person/get"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/handlers/person/list"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/handlers/person/save"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/http/handlers/person/update"
	"github.com/romankravchuk/effective-mobile-test-task/internal/server/middleware"
	"github.com/romankravchuk/effective-mobile-test-task/internal/service/people"
)

func main() {
	cfg, err := config.LoadAPIConfig()
	failedOnError("failed to load config", err)

	log := log.New(cfg.Env, os.Stderr)

	svc, err := people.New(
		people.WithRedisCache(cfg.CacheURL),
		people.WithPostgresPersonStorage(cfg.DatabaseURL),
		people.WithKafkaProducer(&cfg.KafkaMap, cfg.Topic),
	)
	failedOnError("failed to create people service", err)

	mux := chi.NewMux()
	mux.Use(chimiddleware.RequestID)
	mux.Use(middleware.Logger(log))
	mux.Use(chimiddleware.Recoverer)
	mux.Route("/person", func(r chi.Router) {
		r.Post("/", save.New(log, svc))
		r.Get("/", list.New(log, svc))

		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", delete.New(log, svc))
			r.Get("/", get.New(log, svc))
			r.Patch("/", update.New(log, svc))
		})
	})

	log.Info("the api starting", slog.String("port", cfg.Port))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			failedOnError("failed to start the server", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	<-sigCh
	svc.Close()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	failedOnError("failed to shutdown the server", srv.Shutdown(shutdownCtx))

	log.Info("the api stopped")
	os.Exit(0)
}

func failedOnError(msg string, err error) {
	if err != nil {
		slog.Error(msg, sl.Err(err))
		os.Exit(1)
	}
}
