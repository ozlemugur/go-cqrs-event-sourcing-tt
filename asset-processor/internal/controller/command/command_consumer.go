package command

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/consumer"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type CommandConsumer struct {
	reader  *consumer.KafkaConsumer
	handler usecase.CommandHandler
	log     logger.Interface
}

func NewCommandConsumer(reader *consumer.KafkaConsumer, handler usecase.CommandHandler, log logger.Interface) (consumer *CommandConsumer) {
	r := &CommandConsumer{reader, handler, log}
	return r
}

// Start consuming events
func (c *CommandConsumer) Start(ctx context.Context) {
	c.reader.Consume(c.handler.MsgfessageHandler)
}

// Close the Kafka consumer
func (c *CommandConsumer) Close() {
	c.reader.Close()
}
