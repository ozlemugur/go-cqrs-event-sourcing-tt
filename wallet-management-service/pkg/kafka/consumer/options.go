package consumer

import "github.com/confluentinc/confluent-kafka-go/kafka"

type ConsumerOption func(*kafka.ConfigMap) error

// WithAutoOffsetReset configures the auto.offset.reset setting
func WithAutoOffsetReset(offsetReset string) ConsumerOption {
	return func(config *kafka.ConfigMap) error {
		return config.SetKey("auto.offset.reset", offsetReset)
	}
}

// WithMaxPollInterval configures the max.poll.interval.ms setting
func WithMaxPollInterval(interval int) ConsumerOption {
	return func(config *kafka.ConfigMap) error {
		return config.SetKey("max.poll.interval.ms", interval)
	}
}
