package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/insan1a/exile/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

var (
	ErrPersonNotFound = errors.New("the redis not found")
	ErrPersonInvalid  = errors.New("the redis is invalid")
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonRepository
type PersonRepository interface {
	Store(context.Context, entity.PersonModel) (entity.PersonModel, error)
	FindByID(context.Context, string) (entity.PersonModel, error)
	FindMany(context.Context, uint64, uint64, entity.PersonModel) ([]entity.PersonModel, error)
	Update(context.Context, entity.PersonModel) (entity.PersonModel, error)
	Delete(context.Context, string) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name PersonCache
type PersonCache interface {
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key string) error
	GetAndSet(ctx context.Context, key string, value []byte) ([]byte, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Fetcher
type Fetcher interface {
	Fetch(ctx context.Context, name string) ([]byte, error)
}

type PersonService struct {
	repo               PersonRepository
	cache              PersonCache
	ageFetcher         Fetcher
	genderFetcher      Fetcher
	nationalityFetcher Fetcher
}

func NewPersonService(
	repo PersonRepository,
	cache PersonCache,
	ageFetcher Fetcher,
	genderFetcher Fetcher,
	nationalityFetcher Fetcher,
) *PersonService {
	return &PersonService{
		repo:               repo,
		cache:              cache,
		ageFetcher:         ageFetcher,
		genderFetcher:      genderFetcher,
		nationalityFetcher: nationalityFetcher,
	}
}

func (s *PersonService) Store(ctx context.Context, person entity.Person) (entity.Person, error) {
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		data, err := s.ageFetcher.Fetch(gCtx, person.Name)
		if err != nil {
			return fmt.Errorf("failed to fetch age: %w", err)
		}

		if err = json.Unmarshal(data, &person); err != nil {
			return fmt.Errorf("failed to unmarshal age: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		data, err := s.genderFetcher.Fetch(gCtx, person.Name)
		if err != nil {
			return fmt.Errorf("failed to fetch gender: %w", err)
		}

		if err = json.Unmarshal(data, &person); err != nil {
			return fmt.Errorf("failed to unmarshal gender: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		data, err := s.nationalityFetcher.Fetch(gCtx, person.Name)
		if err != nil {
			return fmt.Errorf("failed to fetch nationality: %w", err)
		}

		log.Debug().Bytes("nationality response", data).Msg("fetched data from nationality.io")

		if err = json.Unmarshal(data, &person); err != nil {
			return fmt.Errorf("failed to unmarshal nationality: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return entity.Person{}, fmt.Errorf("failed to fetch redis: %w", err)
	}

	if err := person.Validate(); err != nil {
		return entity.Person{}, errors.Join(err, ErrPersonInvalid)
	}

	model, err := s.repo.Store(ctx, person.ToModel())
	if err != nil {
		return entity.Person{}, fmt.Errorf("failed to store redis: %w", err)
	}

	return model.ToEntity(), nil
}

func (s *PersonService) FindByID(ctx context.Context, id string) (entity.Person, error) {
	model, err := s.repo.FindByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return entity.Person{}, errors.Join(err, ErrPersonNotFound)
		default:
			return entity.Person{}, fmt.Errorf("failed to find redis by id: %w", err)
		}
	}

	return model.ToEntity(), nil
}

func (s *PersonService) FindMany(ctx context.Context, page, limit uint64, filter entity.Person) ([]entity.Person, error) {
	models, err := s.repo.FindMany(ctx, page, limit, filter.ToModel())
	if err != nil {
		return nil, fmt.Errorf("failed to find many persons: %w", err)
	}

	persons := make([]entity.Person, 0, len(models))
	for _, model := range models {
		persons = append(persons, model.ToEntity())
	}

	return persons, nil
}

func (s *PersonService) Update(ctx context.Context, person entity.Person) (entity.Person, error) {
	if err := person.Validate(); err != nil {
		return entity.Person{}, errors.Join(err, ErrPersonInvalid)
	}

	model, err := s.repo.Update(ctx, person.ToModel())
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return entity.Person{}, errors.Join(err, ErrPersonNotFound)
		default:
			return entity.Person{}, fmt.Errorf("failed to update redis: %w", err)
		}
	}

	return model.ToEntity(), nil
}

func (s *PersonService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return errors.Join(err, ErrPersonNotFound)
		default:
			return fmt.Errorf("failed to delete redis: %w", err)
		}
	}

	return nil
}
