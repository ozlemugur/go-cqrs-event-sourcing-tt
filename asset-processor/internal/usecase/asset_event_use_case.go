package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

// AssetUseCase handles asset transactions using event sourcing.
type AssetUseCase struct {
	eventJournal EventJournal // Event producer (Kafka)
	log          logger.Interface
}

// NewAssetUseCase creates a new asset use case.
func NewAssetUseCase(eventJournal EventJournal, l logger.Interface) *AssetUseCase {
	return &AssetUseCase{
		eventJournal: eventJournal,
		log:          l,
	}
}

// Withdraw funds from a wallet (Publishes event)
func (uc *AssetUseCase) Withdraw(ctx context.Context, walletID int, assetName string, amount float64) error {
	event := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  walletID,
		Type:      "withdraw",
		AssetName: assetName,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	if err := uc.eventJournal.PublishEvent(ctx, event); err != nil {
		return fmt.Errorf("Withdraw - PublishEvent: %w", err)
	}

	uc.log.Info("Withdraw event published", "WalletID", walletID, "AssetName", assetName, "Amount", amount)
	return nil
}

// Deposit funds into a wallet (Publishes event)
func (uc *AssetUseCase) Deposit(ctx context.Context, walletID int, assetName string, amount float64) error {
	//wallet check , ya header üzerinden teyitli geldiğini var sayabiliriz ya da httpcall ve ya readonly bir check yapabiliriz
	event := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  walletID,
		Type:      "deposit",
		AssetName: assetName,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	if err := uc.eventJournal.PublishEvent(ctx, event); err != nil {
		return fmt.Errorf("Deposit - PublishEvent: %w", err)
	}

	uc.log.Info("Deposit event published", "WalletID", walletID, "AssetName", assetName, "Amount", amount)
	return nil
}
