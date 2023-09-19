package people

import (
	"context"
	"encoding/json"
	"time"

	kfk "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/validator"
	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
	"github.com/romankravchuk/effective-mobile-test-task/internal/service"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/broker"
	brokerkafka "github.com/romankravchuk/effective-mobile-test-task/internal/storage/broker/kafka"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/cache"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/cache/redis"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person/pg"
)

// Option represents the option for people service
type Option func(s *Service) error

// WithPersonStorage injects user storage into the people service
func WithPersonStorage(people person.Storage) Option {
	return func(s *Service) error {
		if people == nil {
			return service.ErrNilPeopleStorage
		}

		s.people = people
		return nil
	}
}

// WithPostgresPersonStorage injects postgres user storage into the people service
func WithPostgresPersonStorage(url string) Option {
	return func(s *Service) error {
		db, err := storage.NewPostgresPool(url)
		if err != nil {
			return err
		}

		people, err := pg.New(db)
		if err != nil {
			return err
		}

		return WithPersonStorage(people)(s)
	}
}

func WithCache(c cache.Cache, ttl time.Duration) Option {
	return func(s *Service) error {
		s.cache = c
		s.cacheTTL = ttl
		return nil
	}
}

// WithRedisCache injects redis client into the people service
func WithRedisCache(url string) Option {
	return func(s *Service) error {
		client, err := storage.NewRedisClient(url)
		if err != nil {
			return err
		}

		cache, err := redis.New(client)
		if err != nil {
			return err
		}

		return WithCache(cache, 5*time.Minute)(s)
	}
}

func WithProducer(producer broker.Producer, topic string) Option {
	return func(s *Service) error {
		if producer == nil {
			return service.ErrNilProducer
		}

		s.producer = producer
		s.topic = topic
		return nil
	}
}

// WithKafkaProducer injects kafka producer into the people service
func WithKafkaProducer(cfg *kfk.ConfigMap, topic string) Option {
	return func(s *Service) error {
		if cfg == nil {
			return service.ErrNilKafkaConfig
		}

		kafkaProducer, err := storage.NewKafkaProducer(cfg)
		if err != nil {
			return err
		}

		p := brokerkafka.NewProducer(kafkaProducer, topic)
		return WithProducer(p, topic)(s)
	}
}

// Service represents the people service
type Service struct {
	people person.Storage `validate:"required"`

	cache    cache.Cache `validate:"required"`
	cacheTTL time.Duration

	producer broker.Producer
	topic    string
}

// New creates a new people service with given Options.
func New(opts ...Option) (*Service, error) {
	s := &Service{}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	if err := validator.ValidateStruct(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) Save(ctx context.Context, p models.Person) error {
	mp, err := json.Marshal(&p)
	if err != nil {
		return err
	}

	return s.producer.Produce(mp)
}

func (s *Service) Get(ctx context.Context, id string) (*models.Person, error) {
	v, found, err := s.cache.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if found {
		p := models.Person{}
		if err = json.Unmarshal(v, &p); err != nil {
			return nil, err
		}
		return &p, nil
	}

	p, err := s.people.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	mp, _ := json.Marshal(*p)
	if err = s.cache.Set(ctx, id, mp, s.cacheTTL); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) List(ctx context.Context, filter *models.Filter, query string) ([]models.Person, error) {
	v, found, err := s.cache.Get(ctx, query)
	if err != nil {
		return nil, err
	}
	if found {
		p := make([]models.Person, 0)
		if err = json.Unmarshal(v, &p); err != nil {
			return nil, err
		}
		return p, nil
	}

	p, err := s.people.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	mp, _ := json.Marshal(p)
	if err = s.cache.Set(ctx, query, mp, s.cacheTTL); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) Update(ctx context.Context, p *models.Person) error {
	if err := s.people.Update(ctx, p); err != nil {
		return err
	}

	return s.cache.Del(ctx, p.ID)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.people.Delete(ctx, id); err != nil {
		return err
	}

	return s.cache.Del(ctx, id)
}

// Close flushes and closes the producer
func (s *Service) Close() error {
	return s.producer.Close()
}
