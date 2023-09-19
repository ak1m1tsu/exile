package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/romankravchuk/effective-mobile-test-task/internal/config"
	"github.com/romankravchuk/effective-mobile-test-task/internal/lib/sl"
	"github.com/romankravchuk/effective-mobile-test-task/internal/log"
	"github.com/romankravchuk/effective-mobile-test-task/internal/service/person"
	"github.com/romankravchuk/effective-mobile-test-task/internal/storage"
)

func main() {
	cfg, err := config.LoadServiceConfig()
	failedOnError("failed to read config", err)

	log := log.New(cfg.Env, os.Stderr)

	log.Info("config loaded", slog.Any("cfg", cfg))

	consumer, err := storage.NewKafkaConsumer(&cfg.KafkaMap, cfg.Topics)
	failedOnError("failed to create kafka consumer", err)

	producer, err := storage.NewKafkaProducer(&cfg.KafkaMap)
	failedOnError("failed to create kafka producer", err)

	svc, err := person.New(
		person.WithConsumer(consumer),
		person.WithProducer(producer, cfg.Topic),
		person.WithTimeout(cfg.Timeout),
		person.WithPostgresPeopleStorage(cfg.DatabaseURL),
	)
	failedOnError("failed to create service", err)

	log.Info("the service is running")

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, syscall.SIGINT, syscall.SIGTERM)

run:
	for {
		select {
		case <-exitCh:
			producer.Flush(5 * 1000)
			producer.Close()
			consumer.Close()
			log.Info("the service is stopped")
			break run
		default:
			if res, err := svc.Save(); err != nil {
				switch err := err.(type) {
				case kafka.Error:
					if err.Code() == kafka.ErrTimedOut {
						continue
					}
					log.Error("failed to read message", sl.Err(err))
				default:
					if errors.Is(err, person.ErrMessageFromat) || errors.Is(err, person.ErrMessageValidation) {
						if err := svc.SendErrMessage(res, err.Error()); err != nil {
							log.Error("failed to send error message", sl.Err(err))
						}
						log.Info("the invalid message was send", slog.Any("message", res), sl.Err(err))
					} else {
						log.Error("failed to save person", sl.Err(err))
					}
				}
			} else {
				log.Info("the person successfully saved", slog.Any("person", res))
			}
		}
	}
}

func failedOnError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, fmt.Sprintf("error: %v", err))
		os.Exit(1)
	}
}
