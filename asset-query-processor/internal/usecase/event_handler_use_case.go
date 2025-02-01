package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type eventHandler struct {
	repo AssetQueryRepositoryHandler
	log  logger.Interface
}

func NewEventHandler(r AssetQueryRepositoryHandler, l logger.Interface) EventHandler {
	return &eventHandler{repo: r, log: l}
}

// EventTypeHandler defines the function signature for handling events
type EventTypeHandler func(ctx context.Context, repo AssetQueryRepositoryHandler, event entity.WalletEvent) error

// EventHandlers maps event types to their corresponding handler functions
var EventHandlers = map[string]EventTypeHandler{
	"withdraw": handleWithdraw,
	"deposit":  handleDeposit,
	"transfer": handleTransfer,
}

func (h *eventHandler) MsgfessageHandler(key, value []byte) error {
	fmt.Printf("Received message with key: %s, value: %s\n", string(key), string(value))

	ctx := context.Background()

	// Mesajın içeriğine bağlı olarak bir işlem yapın
	if len(value) == 0 {
		return fmt.Errorf("empty message value")
	}

	var event entity.WalletEvent
	if err := json.Unmarshal(value, &event); err != nil {
		h.log.Error(err, "Failed to process event")
	}

	if err := h.ProcessEvent(ctx, event); err != nil {
		h.log.Error(err, "Failed to process event")
	}

	return nil
}

// ProcessEvent dynamically handles wallet events using the map
func (h *eventHandler) ProcessEvent(ctx context.Context, event entity.WalletEvent) error {
	if handler, exists := EventHandlers[event.Type]; exists {
		return handler(ctx, h.repo, event)
	}
	return fmt.Errorf("unknown event type: %s", event.Type)
}

// Withdraw handler
func handleWithdraw(ctx context.Context, repo AssetQueryRepositoryHandler, event entity.WalletEvent) error {
	return repo.UpdateBalance(ctx, event.WalletID, -event.Amount)
}

// Deposit handler
func handleDeposit(ctx context.Context, repo AssetQueryRepositoryHandler, event entity.WalletEvent) error {
	return repo.UpdateBalance(ctx, event.WalletID, event.Amount)
}

// Transfer handler
func handleTransfer(ctx context.Context, repo AssetQueryRepositoryHandler, event entity.WalletEvent) error {
	if err := repo.UpdateBalance(ctx, event.WalletID, -event.Amount); err != nil {
		return err
	}
	return nil //repo.UpdateBalance(ctx, event.TargetWalletID, event.Amount)
}
