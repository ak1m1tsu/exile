package person

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/insan1a/exile/internal/client"
	"github.com/insan1a/exile/internal/lib/validator"
	"github.com/insan1a/exile/internal/models"
	"github.com/insan1a/exile/internal/service"
	"github.com/insan1a/exile/internal/storage"
	"github.com/insan1a/exile/internal/storage/broker"
	"github.com/insan1a/exile/internal/storage/person"
	"github.com/insan1a/exile/internal/storage/person/pg"
	"golang.org/x/sync/errgroup"
)

var (
	ErrMessageFromat     = errors.New("the message have invalid format")
	ErrMessageValidation = errors.New("the message is invalid")
	ErrNilClientFetcher  = errors.New("the client fetcher is nil")
)

type Option func(c *Service) error

func WithTimeout(timeout time.Duration) Option {
	return func(s *Service) error {
		s.timeout = timeout
		return nil
	}
}

func WithConsumer(consumer broker.Consumer) Option {
	return func(c *Service) error {
		if consumer == nil {
			return service.ErrNilConsumer
		}

		c.consumer = consumer
		return nil
	}
}

func WithProducer(producer broker.Producer, topic string) Option {
	return func(s *Service) error {
		if producer == nil {
			return service.ErrNilProducer
		}

		s.producer = producer
		s.producerTopic = topic
		return nil
	}
}

func WithPeopleStorage(people person.Storage) Option {
	return func(s *Service) error {
		if people == nil {
			return service.ErrNilPeopleStorage
		}

		s.people = people
		return nil
	}
}

func WithPostgresPeopleStorage(url string) Option {
	return func(s *Service) error {
		db, err := storage.NewPostgresPool(url)
		if err != nil {
			return err
		}

		storage, err := pg.New(db)
		if err != nil {
			return err
		}

		s.people = storage
		return nil
	}
}

func WithAgifyClient(fetcher client.Fetcher) Option {
	return func(s *Service) error {
		if fetcher == nil {
			return ErrNilClientFetcher
		}

		s.agify = fetcher
		return nil
	}
}

func WithGenderizeClient(fetcher client.Fetcher) Option {
	return func(s *Service) error {
		if fetcher == nil {
			return ErrNilClientFetcher
		}

		s.genderize = fetcher
		return nil
	}
}

func WithNationalizeClient(fetcher client.Fetcher) Option {
	return func(s *Service) error {
		if fetcher == nil {
			return ErrNilClientFetcher
		}

		s.nationalize = fetcher
		return nil
	}
}

type Service struct {
	timeout time.Duration

	consumer      broker.Consumer
	producer      broker.Producer
	producerTopic string

	agify       client.Fetcher
	genderize   client.Fetcher
	nationalize client.Fetcher

	people person.Storage
}

func New(options ...Option) (*Service, error) {
	c := &Service{}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (s *Service) Save(ctx context.Context) ([]byte, error) {
	msg, err := s.consumer.Consume(s.timeout)
	if err != nil {
		return nil, err
	}

	var p models.Person
	if err = json.Unmarshal(msg, &p); err != nil {
		return msg, errors.Join(err, ErrMessageFromat)
	}

	if err = validator.ValidateStruct(p); err != nil {
		return msg, errors.Join(err, ErrMessageValidation)
	}

	clientsCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errs, _ := errgroup.WithContext(clientsCtx)
	errs.Go(func() error {
		data, err := s.nationalize.Fetch(p.Name)
		if err != nil {
			return err
		}

		return json.Unmarshal(data, &p)
	})
	errs.Go(func() error {
		data, err := s.agify.Fetch(p.Name)
		if err != nil {
			return err
		}

		return json.Unmarshal(data, &p)
	})
	errs.Go(func() error {
		data, err := s.genderize.Fetch(p.Name)
		if err != nil {
			return err
		}

		return json.Unmarshal(data, &p)
	})

	if err := errs.Wait(); err != nil {
		return msg, err
	}

	err = s.people.Create(ctx, &p)
	if err != nil {
		return nil, err
	}

	result, _ := json.Marshal(&p)
	return result, nil
}

func (s *Service) SendErrMessage(data []byte, err string) error {
	meta := make(map[string]any)
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}

	errMsg, _ := json.Marshal(&models.ErrorMessage{
		Meta:  meta,
		Error: err,
	})

	return s.producer.Produce(errMsg)
}
