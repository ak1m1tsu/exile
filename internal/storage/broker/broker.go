package broker

import (
	"time"
)

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Consumer --output ./mocks --outpkg mocks
type Consumer interface {
	Consume(timeout time.Duration) ([]byte, error)
	Close() error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name Producer --output ./mocks --outpkg mocks
type Producer interface {
	Produce(message []byte) error
	Close() error
}
