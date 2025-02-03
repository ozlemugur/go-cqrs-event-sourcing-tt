package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/config"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/controller/command"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/usecase"
	eventjournal "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/usecase/event-journal"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/consumer"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/producer"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	kafkaBroker := cfg.Kafka.KAFKA_BROKER // os.Getenv("KAFKA_BROKER")      // e.g., "kafka:9092"
	l.Error("KAFKA_BROKER")
	l.Error(kafkaBroker)
	eventTopic := cfg.Kafka.EVENT_TOPIC // e.g., "event-journal"
	l.Error("EVENT_TOPIC")
	l.Error(eventTopic)

	/*retryTopic := cfg.Kafka.RETRY_TOPIC
	l.Error("RETRY_TOPIC")
	l.Error(retryTopic)

	dlqTopic := cfg.Kafka.DLQ_TOPIC
	l.Error("DLQ_TOPIC")
	l.Error(dlqTopic) */

	// Initialize Kafka Evenjournall*****
	kafkaProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaProducer.Close() // Ensure producer is closed on shutdown

	eventJournal := eventjournal.NewKafkaEventJournal(kafkaProducer, eventTopic)

	assetUseCase := usecase.NewAssetUseCase(
		eventJournal,
		l,
	)

	/*// Initialize RETRY Kafka Producer
	kafkaRetryProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaRetryProducer.Close()

	retryProducer := producer.NewKafkaRetryEventProducer(kafkaRetryProducer, retryTopic)

	// Initialize DLQ Kafka Producer
	kafkaDlqProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaDlqProducer.Close()

	dlqProducer := producer.NewDLQEventProducer(kafkaDlqProducer, dlqTopic) */

	/**********************************************************************************/

	commandTopic := cfg.Kafka.COMMAND_QUEUE_TOPIC
	l.Error("COMMAND-QUEUE")
	l.Error(commandTopic)

	kafkaGroupID := "asset-processor-group" // Consumer group ID
	// Initialize Kafka consumer (command-queue)
	consumer, err := consumer.NewKafkaConsumer(kafkaBroker, kafkaGroupID, commandTopic)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer", "error", err)

	}

	// Initialize Kafka Producer (command-queue)
	commandKafkaProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer commandKafkaProducer.Close() // Ensure producer is closed on shutdown

	commandQueue := command.NewCommandProducer(commandKafkaProducer, commandTopic)

	// Initialize use case (business logic handler)

	commandHandlerUsecase := usecase.NewCommandHandler(commandQueue, assetUseCase, l)

	commandConsumer := command.NewCommandConsumer(consumer, commandHandlerUsecase, l)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer", "error", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	commandConsumer.Start(ctx)

	// Handle system signals for shutdown
	go handleShutdown(cancel, commandConsumer, l)
	// Start consuming events
	l.Info("Starting Event Consumer...")
	commandConsumer.Start(ctx)

}

// handleShutdown gracefully handles shutdown signals.
func handleShutdown(cancel context.CancelFunc, consumer *command.CommandConsumer, log logger.Interface) {
	// Capture system signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Info("Shutdown signal received")

	// Cancel context and close resources
	cancel()
	consumer.Close()

	log.Info("Consumer stopped")
}
