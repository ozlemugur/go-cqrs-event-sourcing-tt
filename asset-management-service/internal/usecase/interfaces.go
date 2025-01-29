package usecase

import (
	"context"
	"time"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/entity"
)

type (
	/* Asset Management UseCase Interface */
	AssetHandler interface {
		Withdraw(ctx context.Context, walletID int, amount float64) error                       // Withdraw funds from a wallet
		Deposit(ctx context.Context, walletID int, amount float64) error                        // Deposit funds into a wallet
		Transfer(ctx context.Context, fromWalletID int, toWalletID int, amount float64) error   // Transfer funds between wallets
		ScheduleTransaction(ctx context.Context, transaction entity.ScheduledTransaction) error // Schedule a future transaction
	}

	/* Asset Repository Interface */
	AssetRepositoryHandler interface {
		GetBalance(ctx context.Context, walletID int) (float64, error)                                              // Fetch wallet balance
		UpdateBalance(ctx context.Context, walletID int, newBalance float64) error                                  // Update wallet balance
		InsertTransaction(ctx context.Context, transaction entity.Transaction) error                                // Store transaction details
		InsertScheduledTransaction(ctx context.Context, transaction entity.ScheduledTransaction) error              // Store a scheduled transaction
		GetScheduledTransactions(ctx context.Context, executeTime time.Time) ([]entity.ScheduledTransaction, error) // Get scheduled transactions that should be executed
		MarkTransactionAsProcessed(ctx context.Context, transactionID int) error                                    // Mark a scheduled transaction as processed
	}

	// EventJournal defines the contract for publishing events.
	EventJournal interface {
		PublishEvent(ctx context.Context, event entity.WalletEvent) error
	}
)
