package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	p     *kafka.Producer
	topic string
}

func NewProducer(p *kafka.Producer, topic string) *Producer {
	return &Producer{
		p:     p,
		topic: topic,
	}
}

func (p *Producer) Produce(msg []byte) error {
	return p.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}, nil)
}

func (p *Producer) Close() error {
	p.p.Flush(1000 * 5)
	p.p.Close()
	return nil
}
