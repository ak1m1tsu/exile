package config

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ilyakaznacheev/cleanenv"
)

type ServiceConfig struct {
	Env string `env:"ENV" env-default:"dev"`

	DatabaseURL string `env:"DATABASE_URL"`

	KafkaMap         kafka.ConfigMap
	GroupID          string        `env:"KAFKA_GROUP_ID"`
	BootstrapServers string        `env:"KAFKA_BOOTSTRAP_SERVERS"`
	AutoOffsetReset  string        `env:"KAFKA_AUTO_OFFSET_RESET"`
	Topic            string        `env:"KAFKA_PRODUCER_TOPIC"`
	Topics           []string      `env:"KAFKA_CONSUMER_TOPICS"`
	Timeout          time.Duration `env:"KAFKA_TIMEOUT" env-default:"100ms"`
}

func LoadServiceConfig() (*ServiceConfig, error) {
	cfg := ServiceConfig{
		KafkaMap: kafka.ConfigMap{},
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	cfg.KafkaMap["group.id"] = cfg.GroupID
	cfg.KafkaMap["auto.offset.reset"] = cfg.AutoOffsetReset
	cfg.KafkaMap["bootstrap.servers"] = cfg.BootstrapServers

	return &cfg, nil
}
