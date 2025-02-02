package usecase

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/entity"
)

type (
	/* Asset Query UseCase Interface */
	AssetHandler interface {
		Withdraw(ctx context.Context, walletID int, amount float64) error                     // Withdraw funds from a wallet
		Deposit(ctx context.Context, walletID int, amount float64) error                      // Deposit funds into a wallet
		Transfer(ctx context.Context, fromWalletID int, toWalletID int, amount float64) error // Transfer funds between wallets
		//ScheduleTransaction(ctx context.Context, transaction entity.ScheduledTransaction) error // Schedule a future transaction
	}

	// AssetQueryRepository defines the interface for querying wallet balances & transactions.
	AssetQueryRepositoryHandler interface {
		// GetBalance retrieves the balance of a specific asset in a wallet
		GetBalance(ctx context.Context, walletID int, assetName string) (float64, error)

		// UpdateBalance updates the balance of a specific asset in a wallet
		UpdateBalance(ctx context.Context, walletID int, assetName string, amount float64) error

		// InsertTransaction inserts a new transaction record into the database
		InsertTransaction(ctx context.Context, txn entity.Transaction) error

		// GetTransactionHistory retrieves the transaction history for a wallet
		GetTransactionHistory(ctx context.Context, walletID int, assetName string) ([]entity.Transaction, error)
	}

	/* Event Handler  UseCase Interface */
	EventHandler interface {
		ProcessEvent(ctx context.Context, event entity.WalletEvent) error
		MsgfessageHandler(key, value []byte) error
	}
)
