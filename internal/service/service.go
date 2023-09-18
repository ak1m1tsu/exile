package service

import "errors"

var (
	ErrNilUserStorage = errors.New("the user storage could not be nil")
	ErrNilKafkaConfig = errors.New("the kafka config could not be nil")
	ErrNilKafkaTopics = errors.New("the kafka topics could not be nil")
	ErrZeroTopics     = errors.New("the kafka topics could not be empty")
)
