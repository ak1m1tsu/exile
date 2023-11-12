package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	rd "github.com/romankravchuk/nix/redis"
)

type PersonCache struct {
	rd *rd.Redis
}

func NewPersonCache(redis *rd.Redis) *PersonCache {
	return &PersonCache{rd: redis}
}

func (c *PersonCache) Get(ctx context.Context, key string) ([]byte, bool, error) {
	res, err := c.rd.Client.Get(ctx, key).Bytes()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, false, nil
		default:
			return nil, false, fmt.Errorf("redis get: %w", err)
		}
	}

	return res, true, nil
}

func (c *PersonCache) Set(ctx context.Context, key string, value []byte) error {
	err := c.rd.Client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

func (c *PersonCache) Delete(ctx context.Context, key string) error {
	err := c.rd.Client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis delete: %w", err)
	}

	return nil
}

func (c *PersonCache) GetAndSet(ctx context.Context, key string, value []byte) ([]byte, error) {
	res, err := c.rd.Client.GetSet(ctx, key, value).Bytes()
	if err != nil {
		return nil, fmt.Errorf("redis getset: %w", err)
	}

	return res, nil
}
