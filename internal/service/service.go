package service

import "errors"

var (
	ErrNilKafkaConfig   = errors.New("the kafka config could not be nil")
	ErrNilKafkaTopics   = errors.New("the kafka topics could not be nil")
	ErrZeroTopics       = errors.New("the kafka topics could not be empty")
	ErrNilConsumer      = errors.New("the kafka consumer could not be nil")
	ErrNilProducer      = errors.New("the kafka producer could not be nil")
	ErrNilPeopleStorage = errors.New("the people storage could not be nil")
)
