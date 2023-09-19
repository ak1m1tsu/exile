package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Consumer struct {
	c *kafka.Consumer
}

func NewConsumer(c *kafka.Consumer) *Consumer {
	return &Consumer{c: c}
}

func (c *Consumer) Consume(timeout time.Duration) ([]byte, error) {
	msg, err := c.c.ReadMessage(timeout)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (c *Consumer) Close() error {
	return c.c.Close()
}
