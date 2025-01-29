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
func (uc *AssetUseCase) Withdraw(ctx context.Context, walletID int, amount float64) error {
	event := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  walletID,
		Type:      "withdraw",
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	if err := uc.eventJournal.PublishEvent(ctx, event); err != nil {
		return fmt.Errorf("Withdraw - PublishEvent: %w", err)
	}

	uc.log.Info("Withdraw event published", "WalletID", walletID, "Amount", amount)
	return nil
}

// Deposit funds into a wallet (Publishes event)
func (uc *AssetUseCase) Deposit(ctx context.Context, walletID int, amount float64) error {
	event := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  walletID,
		Type:      "deposit",
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	if err := uc.eventJournal.PublishEvent(ctx, event); err != nil {
		return fmt.Errorf("Deposit - PublishEvent: %w", err)
	}

	uc.log.Info("Deposit event published", "WalletID", walletID, "Amount", amount)
	return nil
}

// Transfer funds between wallets (Publishes two events: Withdraw + Deposit)
func (uc *AssetUseCase) Transfer(ctx context.Context, fromWalletID, toWalletID int, amount float64) error {
	withdrawEvent := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  fromWalletID,
		Type:      "withdraw",
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	depositEvent := entity.WalletEvent{
		EventID:   uuid.New().String(),
		WalletID:  toWalletID,
		Type:      "deposit",
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

	uc.log.Info("Transfer events published", "FromWalletID", fromWalletID, "ToWalletID", toWalletID, "Amount", amount)
	return nil
}

// Schedule a future transaction (Publishes event)
func (uc *AssetUseCase) ScheduleTransaction(ctx context.Context, transaction entity.ScheduledTransaction) error {
	//transaction.ID = uuid.New().String()
	/*	event := entity.ScheduledTransactionEvent{
			EventID:     transaction.ID,
			WalletID:    transaction.FromWallet,
			Type:        ,
			Amount:      transaction.Amount,
			ExecuteTime: transaction.ExecuteTime.Unix(),
		}

		if err := uc.eventJournal.PublishEvent(ctx, event); err != nil {
			return fmt.Errorf("ScheduleTransaction - PublishEvent: %w", err)
		}
	*/
	uc.log.Info("Scheduled transaction event published", "TransactionID", transaction.ID, "ExecuteTime", transaction.ExecuteTime)
	return nil
}

/*package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type AssetUseCase struct {
	repo AssetRepositoryHandler
	log  logger.Interface
}

// NewAssetUseCase - Creates a new asset use case
func NewAssetUseCase(r AssetRepositoryHandler, l logger.Interface) *AssetUseCase {
	return &AssetUseCase{
		repo: r,
		log:  l,
	}
}


// Withdraw funds from a wallet
func (uc *AssetUseCase) Withdraw(ctx context.Context, walletID int, amount float64) error {
	balance, err := uc.repo.GetBalance(ctx, walletID)
	if err != nil {
		return fmt.Errorf("Withdraw - GetBalance: %w", err)
	}

	if balance < amount {
		return fmt.Errorf("Withdraw - insufficient funds")
	}

	newBalance := balance - amount
	if err := uc.repo.UpdateBalance(ctx, walletID, newBalance); err != nil {
		return fmt.Errorf("Withdraw - UpdateBalance: %w", err)
	}

	transaction := entity.Transaction{
		WalletID:  walletID,
		Type:      "withdraw",
		Amount:    amount,
		Status:    "completed",
		CreatedAt: time.Now(),
	}
	if err := uc.repo.InsertTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("Withdraw - InsertTransaction: %w", err)
	}

	uc.log.Info("Withdraw successful", "WalletID", walletID, "Amount", amount)
	return nil
}

// Deposit funds into a wallet
func (uc *AssetUseCase) Deposit(ctx context.Context, walletID int, amount float64) error {
	balance, err := uc.repo.GetBalance(ctx, walletID)
	if err != nil {
		return fmt.Errorf("Deposit - GetBalance: %w", err)
	}

	newBalance := balance + amount
	if err := uc.repo.UpdateBalance(ctx, walletID, newBalance); err != nil {
		return fmt.Errorf("Deposit - UpdateBalance: %w", err)
	}

	transaction := entity.Transaction{
		WalletID:  walletID,
		Type:      "deposit",
		Amount:    amount,
		Status:    "completed",
		CreatedAt: time.Now(),
	}
	if err := uc.repo.InsertTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("Deposit - InsertTransaction: %w", err)
	}

	uc.log.Info("Deposit successful", "WalletID", walletID, "Amount", amount)
	return nil
}

// Transfer funds between wallets
func (uc *AssetUseCase) Transfer(ctx context.Context, fromWalletID, toWalletID int, amount float64) error {
	// Withdraw from sender wallet
	if err := uc.Withdraw(ctx, fromWalletID, amount); err != nil {
		return fmt.Errorf("Transfer - Withdraw: %w", err)
	}

	// Deposit into receiver wallet
	if err := uc.Deposit(ctx, toWalletID, amount); err != nil {
		// Rollback: Revert the withdrawal if deposit fails
		_ = uc.Deposit(ctx, fromWalletID, amount)
		return fmt.Errorf("Transfer - Deposit: %w", err)
	}

	uc.log.Info("Transfer successful", "FromWalletID", fromWalletID, "ToWalletID", toWalletID, "Amount", amount)
	return nil
}

// Schedule a future transaction
func (uc *AssetUseCase) ScheduleTransaction(ctx context.Context, transaction entity.ScheduledTransaction) error {
	if err := uc.repo.InsertScheduledTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("ScheduleTransaction - InsertScheduledTransaction: %w", err)
	}

	uc.log.Info("Scheduled transaction created", "TransactionID", transaction.ID, "ExecuteTime", transaction.ExecuteTime)
	return nil
}
*/
