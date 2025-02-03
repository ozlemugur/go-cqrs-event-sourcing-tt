package consumer

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	reader *kafka.Consumer
}

// NewKafkaConsumer initializes a new Kafka consumer
func NewKafkaConsumer(broker, groupID string, topic string, opts ...ConsumerOption) (*KafkaConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest", // Default offset reset policy
	}

	// Apply consumer options
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, fmt.Errorf("failed to apply consumer option: %w", err)
		}
	}

	// Create the Kafka consumer
	reader, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	// Subscribe to topic
	if err := reader.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &KafkaConsumer{reader: reader}, nil
}

// Consume listens for messages and processes them with the given handler
func (c *KafkaConsumer) Consume(handler func(key, value []byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(-1)
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}
		if err := handler(msg.Key, msg.Value); err != nil {
			return fmt.Errorf("handler error: %w", err)
		}
	}
}

// Close closes the Kafka consumer
func (c *KafkaConsumer) Close() {
	c.reader.Close()
}
