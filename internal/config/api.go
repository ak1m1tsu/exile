package config

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ilyakaznacheev/cleanenv"
)

type APIConfig struct {
	Env string `env:"ENV" env-default:"dev"`

	Port         string        `env:"PORT" env-default:"8080"`
	IdleTimeout  time.Duration `env:"API_IDLE_TIMEOUT" env-default:"30s"`
	ReadTimeout  time.Duration `env:"API_READ_TIMEOUT" env-default:"5s"`
	WriteTimeout time.Duration `env:"API_WRITE_TIMEOUT" env-default:"5s"`

	DatabaseURL string `env:"DATABASE_URL"`
	CacheURL    string `env:"CACHE_URL"`

	KafkaMap         kafka.ConfigMap
	BootstrapServers string `env:"KAFKA_BOOTSTRAP_SERVERS"`
	Topic            string `env:"KAFKA_PRODUCER_TOPIC"`
}

func LoadAPIConfig() (*APIConfig, error) {
	path, err := getPath("API_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	cfg := APIConfig{
		KafkaMap: kafka.ConfigMap{},
	}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	cfg.KafkaMap["bootstrap.servers"] = cfg.BootstrapServers

	return &cfg, nil
}
