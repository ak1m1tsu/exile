package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/romankravchuk/effective-mobile-test-task/internal/config"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage"
)

type Person struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func main() {
	cfg, err := config.LoadAPIConfig()
	failedOnError(err, "failed to read config file")

	topic := cfg.Topic
	p, err := storage.NewKafkaProducer(&cfg.KafkaMap)
	failedOnError(err, "failed to create a kafka producer")

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					slog.Error(
						"failed to deliver message",
						slog.Any("tp", ev.TopicPartition),
					)
				} else {
					slog.Info(
						"produced event to topic",
						slog.Group(
							"topic",
							slog.String("name", *ev.TopicPartition.Topic),
							slog.String("key", string(ev.Key)),
							slog.String("value", string(ev.Value)),
						),
					)
				}
			}
		}
	}()

	var wg sync.WaitGroup
	for n := 0; n < 10; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := uuid.New()
			person := Person{
				Name:    faker.FirstName(),
				Surname: faker.LastName(),
			}

			marshaledPerson, err := json.Marshal(&person)
			if err != nil {
				slog.Error(
					"failed to marshal person",
					slog.Any("person", person),
					slog.String("error", err.Error()),
				)
				return
			}

			err = p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topic,
					Partition: kafka.PartitionAny,
				},
				Key:   []byte(key.String()),
				Value: marshaledPerson,
			}, nil)
			if err != nil {
				slog.Error(
					"failed to produce message to kafka",
					slog.String("error", err.Error()),
				)
			}
		}()
	}

	wg.Wait()
	p.Flush(15 * 1000)
	p.Close()
	os.Exit(0)
}

func failedOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, slog.String("error", err.Error()))
		os.Exit(1)
	}
}
