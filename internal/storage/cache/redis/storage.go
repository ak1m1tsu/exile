package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage"
)

type Storage struct {
	client *redis.Client
}

func New(client *redis.Client) (*Storage, error) {
	if client == nil {
		return nil, storage.ErrNilRedisClient
	}

	return &Storage{client: client}, nil
}

func (s *Storage) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if err := s.client.Set(ctx, key, value, ttl).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}

func (s *Storage) Get(ctx context.Context, key string) ([]byte, bool, error) {
	res, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return res, true, nil
}

func (s *Storage) Del(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}
	return nil
}
