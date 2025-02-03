package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/producer"
)

// CommandProducer handles publishing commands to Kafka.
type CommandProducer struct {
	producer *producer.KafkaProducer
	topic    string
}

// NewCommandProducer creates a new Kafka-based command producer.
func NewCommandProducer(kafkaProducer *producer.KafkaProducer, topic string) *CommandProducer {
	return &CommandProducer{
		producer: kafkaProducer,
		topic:    topic,
	}
}

// PublishCommand serializes and sends a command to Kafka.
func (c *CommandProducer) PublishCommand(ctx context.Context, command interface{}) error {
	commandBytes, err := json.Marshal(command)
	if err != nil {
		return fmt.Errorf("failed to serialize command: %w", err)
	}

	err = c.producer.ProduceEvent(c.topic, commandBytes, fmt.Sprintf("command-%s", command))
	if err != nil {
		log.Printf("Failed to publish command: %v", err)
		return err
	}

	log.Printf("Command published: %+v", command)
	return nil
}
