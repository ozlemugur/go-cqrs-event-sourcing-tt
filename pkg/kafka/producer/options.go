package producer

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ProducerOption func(*kafka.ConfigMap) error

// WithAcks configures the producer's acks
func WithAcks(acks string) ProducerOption {
	return func(config *kafka.ConfigMap) error {
		return config.SetKey("acks", acks)
	}
}

// WithCompression configures the producer's compression type
func WithCompression(compression string) ProducerOption {
	return func(config *kafka.ConfigMap) error {
		return config.SetKey("compression.type", compression)
	}
}
