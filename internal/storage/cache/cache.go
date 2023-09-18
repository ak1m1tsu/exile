package cache

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("the cached value not found")

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Cache --output ./mocks --outpkg mocks
type Cache interface {
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Get(ctx context.Context, key string) ([]byte, bool, error)
	Del(ctx context.Context, key string) error
}
