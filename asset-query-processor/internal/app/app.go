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

	// Initialize use case (business logic handler)
	queryRepo := repo.NewAssetQueryRepo(pg)
	eventHandler := usecase.NewEventHandler(queryRepo, l)

	// Kafka setup
	kafkaBroker := os.Getenv("KAFKA_BROKER")      // e.g., "kafka:9092"
	eventTopic := os.Getenv("EVENT_TOPIC")        // e.g., "event-journal"
	kafkaGroupID := "asset-query-processor-group" // Consumer group ID

	// Initialize Kafka consumer
	consumer, err := controller.NewEventConsumer(kafkaBroker, kafkaGroupID, eventTopic, eventHandler, l)
	if err != nil {
		log.Fatal("Failed to initialize Kafka consumer", "error", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle system signals for shutdown
	go handleShutdown(cancel, consumer, l)

	// Start consuming events
	l.Info("Starting Event Consumer...")
	consumer.Start(ctx)

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
