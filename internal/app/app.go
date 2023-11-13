package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/insan1a/exile/internal/adapter/api"
	rd "github.com/insan1a/exile/internal/adapter/cache/redis"
	pg "github.com/insan1a/exile/internal/adapter/database/postgres"
	"github.com/insan1a/exile/internal/controller/http/response"
	v1 "github.com/insan1a/exile/internal/controller/http/v1"
	"github.com/insan1a/exile/internal/domain/service"
	"github.com/romankravchuk/nix/httpserver"
	log "github.com/romankravchuk/nix/log/zerolog"
	"github.com/romankravchuk/nix/postgres"
	"github.com/romankravchuk/nix/redis"
)

func Run(cfg *Config) error {
	logger := log.New(os.Stdout, log.Info)

	db, err := postgres.New(cfg.Database.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	cache, err := redis.New(cfg.Cache.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to cache: %w", err)
	}

	ageFetcher := api.NewAgeFetcher(cfg.API.Age.URL)
	genderFetcher := api.NewGenderFetcher(cfg.API.Gender.URL)
	nationalityFetcher := api.NewNationalityFetcher(cfg.API.Nationality.URL)

	personRepository := pg.NewPersonRepository(db)
	personCache := rd.NewPersonCache(cache)

	personService := service.NewPersonService(
		personRepository,
		personCache,
		ageFetcher,
		genderFetcher,
		nationalityFetcher,
	)

	mux := chi.NewMux()
	mux.NotFound(response.NotFound)
	mux.MethodNotAllowed(response.MethodNotAllowed)
	v1.MountPersonHandler(mux, personService, logger)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	server := httpserver.New(ctx, mux, httpserver.Port(cfg.HTTP.Port))

	return server.Run(ctx)
}
