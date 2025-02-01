package controller

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/consumer"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type EventConsumer struct {
	reader  *consumer.KafkaConsumer
	handler usecase.EventHandler
	log     logger.Interface
}

func NewEventConsumer(reader *consumer.KafkaConsumer, handler usecase.EventHandler, log logger.Interface) (consumer *EventConsumer) {
	r := &EventConsumer{reader, handler, log}
	return r
}

// Start consuming events
func (c *EventConsumer) Start(ctx context.Context) {
	c.reader.Consume(c.handler.MsgfessageHandler)
}

// Close the Kafka consumer
func (c *EventConsumer) Close() {
	c.reader.Close()
}
