package main

import (
	"log"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/config"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
