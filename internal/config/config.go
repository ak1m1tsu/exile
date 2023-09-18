package config

import (
	"errors"
	"os"
)

var ErrEmptyPath = errors.New("the config path is empty")

func getPath(env string) (string, error) {
	path := os.Getenv(env)
	if path == "" {
		return "", ErrEmptyPath
	}

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}
