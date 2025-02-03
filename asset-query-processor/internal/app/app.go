package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/config"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/controller"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/usecase/repo"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/consumer"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/kafka/producer"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/postgres"
)

// Asset Query Processor (Consumes Kafka Events and Updates Query DB)
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	l.Error("postgresql ayakta.")

	// Kafka setup

	kafkaBroker := cfg.Kafka.KAFKA_BROKER // os.Getenv("KAFKA_BROKER")      // e.g., "kafka:9092"
	l.Error("KAFKA_BROKER")
	l.Error(kafkaBroker)
	eventTopic := cfg.Kafka.EVENT_TOPIC // e.g., "event-journal"

	l.Error("EVENT_TOPIC")
	l.Error(eventTopic)
	kafkaGroupID := "asset-query-processor-group" // Consumer group ID

	retryTopic := cfg.Kafka.RETRY_TOPIC
	l.Error("RETRY_TOPIC")
	l.Error(retryTopic)

	dlqTopic := cfg.Kafka.DLQ_TOPIC
	l.Error("DLQ_TOPIC")
	l.Error(dlqTopic)

	// Initialize RETRY Kafka Producer
	kafkaRetryProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaRetryProducer.Close()

	retryProducer := controller.NewKafkaRetryEventProducer(kafkaRetryProducer, retryTopic)

	// Initialize DLQ Kafka Producer
	kafkaDlqProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaDlqProducer.Close()

	dlqProducer := controller.NewDLQEventProducer(kafkaDlqProducer, dlqTopic)

	/**********************************************************************************/

	// Initialize Kafka consumer
	consumer, err := consumer.NewKafkaConsumer(kafkaBroker, kafkaGroupID, eventTopic)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer", "error", err)

	}

	// Initialize use case (business logic handler)
	queryRepo := repo.NewAssetQueryRepo(pg)
	eventHandler := usecase.NewEventHandler(queryRepo, retryProducer, dlqProducer, l)

	eventConsumer := controller.NewEventConsumer(consumer, eventHandler, l)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer", "error", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventConsumer.Start(ctx)

	// Handle system signals for shutdown
	go handleShutdown(cancel, eventConsumer, l)
	// Start consuming events
	l.Info("Starting Event Consumer...")
	eventConsumer.Start(ctx)

}

// handleShutdown gracefully handles shutdown signals.
func handleShutdown(cancel context.CancelFunc, consumer *controller.EventConsumer, log logger.Interface) {
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
