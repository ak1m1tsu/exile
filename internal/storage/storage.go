package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/errtools"
)

var (
	ErrNilDB                = errors.New("the database pool is nil")
	ErrURLEmpty             = errors.New("the url is empty")
	ErrUnsupportedParamType = errors.New("the query param have unsupported type")
)

// NewPostgresPool creates a new database connection pool for PostgreSQL.
//
// If url is empty ErrURLEmpty is returned.
func NewPostgresPool(url string) (*sql.DB, error) {
	const op = "storage.NewPostgresPool"

	if url == "" {
		return nil, errtools.WithOperation(ErrURLEmpty, op)
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return db, nil
}

// NewRedisClient creates a new Redis client.
//
// If url is empty ErrURLEmpty is returned.
func NewRedisClient(url string) (*redis.Client, error) {
	const op = "storage.NewRedisClient"

	if url == "" {
		return nil, errtools.WithOperation(ErrURLEmpty, op)
	}

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = client.Ping(ctx).Err(); err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return client, nil
}

// NewKafkaConsumer creates a new kafka consumer.
//
// Consumer subscribes to the given topics.
func NewKafkaConsumer(cfg *kafka.ConfigMap, topics []string) (*kafka.Consumer, error) {
	const op = "storage.NewKafkaConsumer"

	consumer, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	if err = consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return consumer, nil
}

// NewKafkaProducer creates a new kafka producer.
func NewKafkaProducer(cfg *kafka.ConfigMap) (*kafka.Producer, error) {
	const op = "storage.NewKafkaProducer"

	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, errtools.WithOperation(err, op)
	}

	return producer, nil
}
