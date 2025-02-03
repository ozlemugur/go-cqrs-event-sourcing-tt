package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/config"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/controller/command"
	v1 "github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/controller/http/v1"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/usecase"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/httpserver"
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

	// Initialize Kafka Producer
	kafkaProducer, err := producer.NewKafkaProducer(kafkaBroker)
	if err != nil {
		l.Fatal(" Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaProducer.Close() // Ensure producer is closed on shutdown

	commandQueue := command.NewCommandProducer(kafkaProducer, eventTopic)

	assetUseCase := usecase.NewAssetUseCase(
		commandQueue,
		l,
	)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, assetUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:

		l.Info("app - Run - signal: " + s.String())

	case err = <-httpServer.Notify():

		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()

	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
