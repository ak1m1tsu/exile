package app

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP struct {
		Port         string        `yaml:"port" env:"PORT" env-default:"8080"`
		IdleTimeout  time.Duration `yaml:"idleTimeout" env:"IDLE_TIMEOUT" env-default:"60s"`
		ReadTimeout  time.Duration `yaml:"readTimeout" env:"READ_TIMEOUT" env-default:"10s"`
		WriteTimeout time.Duration `yaml:"writeTimeout" env:"WRITE_TIMEOUT" env-default:"10s"`
	} `yaml:"http"`

	Database struct {
		URL string `yaml:"url" env:"DATABASE_URL"`
	} `yaml:"database"`

	Cache struct {
		URL string `yaml:"url" env:"CACHE_URL"`
	} `yaml:"cache"`

	API struct {
		Age struct {
			URL string `yaml:"url" env:"AGE_URL"`
		} `yaml:"age"`
		Gender struct {
			URL string `yaml:"url" env:"GENDER_URL"`
		}
		Nationality struct {
			URL string `yaml:"url" env:"NATIONALITY_URL"`
		}
	} `yaml:"api"`

	Kafka struct {
		GroupID          string `yaml:"groupId" env:"KAFKA_GROUP_ID"`
		BootstrapServers string `yaml:"bootstrapServers" env:"KAFKA_BOOTSTRAP_SERVERS"`
		Consumer         struct {
			Topics []string `yaml:"topics" env:"KAFKA_CONSUMER_TOPICS"`
		} `yaml:"consumer"`
		Producer struct {
			Topic string `yaml:"topic" env:"KAFKA_PRODUCER_TOPIC"`
		} `yaml:"producer"`
	} `yaml:"kafka"`
}

var ErrConfigPath = errors.New("config path is empty")

func LoadConfig(path string) (*Config, error) {
	if path == "" {
		return nil, ErrConfigPath
	}

	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("config file not found: %w", err)
	}

	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to read env: %w", err)
	}

	return cfg, nil
}
