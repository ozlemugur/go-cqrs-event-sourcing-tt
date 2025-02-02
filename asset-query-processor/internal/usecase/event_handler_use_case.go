package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type eventHandler struct {
	repo  AssetQueryRepositoryHandler
	log   logger.Interface
	retry RetryEventProducer
	dlq   DLQEventProducer
}

func NewEventHandler(r AssetQueryRepositoryHandler, retry RetryEventProducer, dlq DLQEventProducer, l logger.Interface) EventHandler {
	return &eventHandler{repo: r, log: l, retry: retry, dlq: dlq}
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

	if len(value) == 0 {
		return fmt.Errorf("empty message value")
	}

	// Base64 kodlu mesajı çözme
	decodedValue, err := base64.StdEncoding.DecodeString(strings.Trim(string(value), "\""))
	if err != nil {
		h.log.Error(err, "Base64 decode error")
		return h.retryOrSendToDLQ(ctx, key, value, "Base64 decode error")
	}

	// JSON mesajını çözme
	var event entity.WalletEvent
	if err := json.Unmarshal(decodedValue, &event); err != nil {
		h.log.Error(err, "JSON unmarshal error")
		return h.retryOrSendToDLQ(ctx, key, decodedValue, "JSON unmarshal error")
	}

	// Event'i işleme
	if err := h.ProcessEvent(ctx, event); err != nil {
		h.log.Error(err, "Failed to process event")
		return h.retryOrSendToDLQ(ctx, key, decodedValue, "Event processing error")
	}

	return nil
}

func (h *eventHandler) retryOrSendToDLQ(ctx context.Context, key, value []byte, errorMsg string) error {
	const maxRetries = 3

	// Retry deneme sayısını belirle (örnek olarak 0 başlatıldı)
	retryCount := h.getRetryCount(key)

	if retryCount < maxRetries {
		h.log.Info("Retrying message", "key", string(key), "retryCount", retryCount)
		h.incrementRetryCount(key)
		return h.retry.PublishRetryEvent(ctx, value)
	} else {
		h.log.Info("Sending message to DLQ", "key", string(key))
		return h.dlq.PublishDLQEvent(ctx, value)
	}
}

func (h *eventHandler) getRetryCount(key []byte) int {
	return 3
}

func (h *eventHandler) incrementRetryCount(key []byte) {
	// Retry sayısını bir metadata olarak artır
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
	return repo.UpdateBalance(ctx, event.WalletID, event.AssetName, -event.Amount)
}

// Deposit handler
func handleDeposit(ctx context.Context, repo AssetQueryRepositoryHandler, event entity.WalletEvent) error {
	return repo.UpdateBalance(ctx, event.WalletID, event.AssetName, event.Amount)
}

// Transfer handler
func handleTransfer(ctx context.Context, repo AssetQueryRepositoryHandler, event entity.WalletEvent) error {
	// Withdraw from sender wallet
	if err := repo.UpdateBalance(ctx, event.WalletID, event.AssetName, -event.Amount); err != nil {
		return err
	}
	// Deposit to target wallet (Assume another field such as TargetWalletID exists for a transfer event)
	// if err := repo.UpdateBalance(ctx, event.TargetWalletID, event.AssetName, event.Amount); err != nil {
	//     return err
	// }
	return nil
}
