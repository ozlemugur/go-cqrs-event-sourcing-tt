package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaProducer represents a Kafka producer for event sourcing
type KafkaProducer struct {
	writer *kafka.Producer
}

// NewKafkaProducer initializes a new Kafka producer
func NewKafkaProducer(broker string, opts ...ProducerOption) (*KafkaProducer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":  broker,
		"acks":               "all",    // Ensure message durability
		"compression.type":   "snappy", // Optimize performance
		"retries":            5,        // Retry on failure
		"retry.backoff.ms":   500,      // Delay between retries
		"enable.idempotence": true,     // Ensure exactly-once delivery
	}

	// Apply producer options
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, fmt.Errorf("failed to apply producer option: %w", err)
		}
	}

	// Create the Kafka producer
	writer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{writer: writer}, nil
}

// ProduceEvent sends a structured event to Kafka
func (p *KafkaProducer) ProduceEvent(topic string, event interface{}, key string, headers []kafka.Header) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key), // Partition by key (e.g., wallet ID)
		Value:          eventBytes,
		Headers:        headers, // Header bilgilerini ekle
	}

	// Deliver event asynchronously
	go func() {
		err := p.writer.Produce(msg, nil)
		if err != nil {
			log.Printf("Failed to produce event: %v", err)
		}
	}()

	return nil
}

// Close closes the Kafka producer
func (p *KafkaProducer) Close() {
	p.writer.Close()
}
