package person

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/romankravchuk/effective-mobile-test-task/internal/client"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/validator"
	"github.com/romankravchuk/effective-mobile-test-task/internal/models"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage/person/pg"
	"golang.org/x/sync/errgroup"
)

var (
	ErrNilConsumer       = errors.New("the kafka consumer could not be nil")
	ErrNilProducer       = errors.New("the kafka producer could not be nil")
	ErrNilPeopleStorage  = errors.New("the people storage could not be nil")
	ErrMessageFromat     = errors.New("the message have invalid format")
	ErrMessageValidation = errors.New("the message is invalid")
)

type Option func(c *Service) error

func WithTimeout(timeout time.Duration) Option {
	return func(s *Service) error {
		s.timeout = timeout
		return nil
	}
}

func WithConsumer(consumer *kafka.Consumer) Option {
	return func(c *Service) error {
		if consumer == nil {
			return ErrNilConsumer
		}

		c.consumer = consumer
		return nil
	}
}

func WithProducer(producer *kafka.Producer, topic string) Option {
	return func(s *Service) error {
		if producer == nil {
			return ErrNilProducer
		}

		s.producer = producer
		s.producerTopic = topic
		return nil
	}
}

func WithPeopleStorage(people person.Storage) Option {
	return func(s *Service) error {
		if people == nil {
			return ErrNilPeopleStorage
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

type Service struct {
	timeout time.Duration

	consumer      *kafka.Consumer
	producer      *kafka.Producer
	producerTopic string

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

func (s *Service) Save() ([]byte, error) {
	msg, err := s.consumer.ReadMessage(s.timeout)
	if err != nil {
		return nil, err
	}

	var p models.Person
	if err = json.Unmarshal(msg.Value, &p); err != nil {
		return msg.Value, errors.Join(err, ErrMessageFromat)
	}

	if err = validator.ValidateStruct(p); err != nil {
		return msg.Value, errors.Join(err, ErrMessageValidation)
	}

	clientsCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errs, _ := errgroup.WithContext(clientsCtx)
	errs.Go(func() error {
		p.Nationality, err = client.FetchNationality(p.Name)
		return err
	})
	errs.Go(func() error {
		age, err := client.FetchAge(p.Name)
		p.Age = age
		return err
	})
	errs.Go(func() error {
		p.Gender, err = client.FetchGender(p.Name)
		return err
	})

	if err := errs.Wait(); err != nil {
		return msg.Value, err
	}

	createCtx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	err = s.people.Create(createCtx, p)
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

	return s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &s.producerTopic,
			Partition: kafka.PartitionAny,
		},
		Value: errMsg,
	}, nil)
}
