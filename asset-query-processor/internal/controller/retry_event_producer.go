package controller

import (
	"context"
	"log"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/producer"
)

// KafkaEventProducer is the implementation of EventJournal using Kafka.
type KafkaRetryEventProducer struct {
	producer *producer.KafkaProducer
	topic    string
}

// KafkaEventProducer creates a new Kafka-based event journal.
func NewKafkaRetryEventProducer(kafkaProducer *producer.KafkaProducer, topic string) *KafkaRetryEventProducer {
	return &KafkaRetryEventProducer{
		producer: kafkaProducer,
		topic:    topic,
	}
}

// PublishEvent serializes and sends the event to Kafka.
/*func (e *KafkaRetryEventProducer) PublishRetryEvent(ctx context.Context, event entity.WalletEvent) error {
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

func (e *KafkaRetryEventProducer) PublishRetryEvent(ctx context.Context, event []byte) error {
	err := e.producer.ProduceEvent(e.topic, event, "")
	if err != nil {
		log.Printf("Failed to publish retry event: %v", err)
		return err
	}

	log.Println("Retry event published")
	return nil
}
