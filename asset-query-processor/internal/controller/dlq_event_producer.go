package controller

import (
	"context"
	"log"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/producer"
)

type DLQEventProducer struct {
	producer *producer.KafkaProducer
	topic    string
}

// KafkaEventProducer creates a new Kafka-based event journal.
func NewDLQEventProducer(kafkaProducer *producer.KafkaProducer, topic string) *DLQEventProducer {
	return &DLQEventProducer{
		producer: kafkaProducer,
		topic:    topic,
	}
}

/*
// PublishEvent serializes and sends the event to Kafka.
func (e *DLQEventProducer) PublishDLQEvent(ctx context.Context, event entity.WalletEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	err = e.producer.ProduceEvent(e.topic, eventBytes, fmt.Sprintf("wallet-%d", event.WalletID))
	if err != nil {
		log.Printf("Failed to publish event: %v", err)
		return err
	}

	log.Printf("Event published: %+v", event)
	return nil
} */

func (e *DLQEventProducer) PublishDLQEvent(ctx context.Context, event []byte) error {
	err := e.producer.ProduceEvent(e.topic, event, "")
	if err != nil {
		log.Printf("Failed to publish DLQ event: %v", err)
		return err
	}

	log.Println("DLQ event published")
	return nil
}
