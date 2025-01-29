package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type EventConsumer struct {
	reader  *kafka.Consumer
	handler usecase.EventHandler
	log     logger.Interface
}

func NewEventConsumer(broker, groupID, topic string, handler usecase.EventHandler, log logger.Interface) (*EventConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	reader, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	if err := reader.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	return &EventConsumer{reader: reader, handler: handler, log: log}, nil
}

// Start consuming events
func (c *EventConsumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(-1)
		if err != nil {
			c.log.Error(err, "Failed to read Kafka message")
			continue
		}

		var event entity.WalletEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			c.log.Error(err, "Failed to deserialize event")
			continue
		}

		c.log.Info("Processing event", "event_id", event.EventID, "type", event.Type)
		if err := c.handler.ProcessEvent(ctx, event); err != nil {
			c.log.Error(err, "Failed to process event")
		}
	}
}

// Close the Kafka consumer
func (c *EventConsumer) Close() {
	c.reader.Close()
}
