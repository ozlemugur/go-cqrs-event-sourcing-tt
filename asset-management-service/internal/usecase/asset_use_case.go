package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/entity"
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

// Transfer funds between wallets (Publishes two events: Withdraw + Deposit)
func (uc *AssetUseCase) Transfer(ctx context.Context, fromWalletID, toWalletID int, assetName string, amount float64) error {
	// TODO: balance check gerekiyor
	withdrawEvent := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  fromWalletID,
		Type:      "withdraw",
		AssetName: assetName,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	depositEvent := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  toWalletID,
		Type:      "deposit",
		AssetName: assetName,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	// Publish both events
	if err := uc.eventJournal.PublishEvent(ctx, withdrawEvent); err != nil {
		return fmt.Errorf("Transfer - Withdraw PublishEvent: %w", err)
	}
	if err := uc.eventJournal.PublishEvent(ctx, depositEvent); err != nil {
		return fmt.Errorf("Transfer - Deposit PublishEvent: %w", err)
	}

	uc.log.Info("Transfer events published", "FromWalletID", fromWalletID, "ToWalletID", toWalletID, "AssetName", assetName, "Amount", amount)
	return nil
}

/*
// Schedule a future transaction (Publishes event) this feature will implement differently
func (uc *AssetUseCase) ScheduleTransaction(ctx context.Context, transaction entity.ScheduledTransactionRequest) error {
	//transaction.ID = uuid.New().String()

	/*event := entity.ScheduledTransactionEvent{
		EventID:     transaction.ID,
		WalletID:    transaction.FromWallet,
		Type:        "scheduled_transaction",
		AssetName:   transaction.AssetName,
		Amount:      transaction.Amount,
		ExecuteTime: transaction.ExecuteTime.Unix(),
	} */

/*if err := uc.eventJournal.PublishEvent(ctx, event); err != nil {
	return fmt.Errorf("ScheduleTransaction - PublishEvent: %w", err)
} */

//uc.log.Info("Scheduled transaction event published", "TransactionID", transaction.ID, "AssetName", transaction.AssetName, "ExecuteTime", transaction.ExecuteTime)
//} */
