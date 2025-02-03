package eventjournal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/producer"
)

// KafkaEventJournal is the implementation of EventJournal using Kafka.
type KafkaEventJournal struct {
	producer *producer.KafkaProducer
	topic    string
}

// NewKafkaEventJournal creates a new Kafka-based event journal.
func NewKafkaEventJournal(kafkaProducer *producer.KafkaProducer, topic string) *KafkaEventJournal {
	return &KafkaEventJournal{
		producer: kafkaProducer,
		topic:    topic,
	}
}

// PublishEvent serializes and sends the event to Kafka.
func (e *KafkaEventJournal) PublishEvent(ctx context.Context, event entity.WalletEvent) error {
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
}
