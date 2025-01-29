// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/entity"
)

type (
	/* Wallet Query UseCase Interface */
	WalletQueryHandler interface {
		GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
		GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error)
		// Balance & Transactions
		GetBalance(ctx context.Context, walletID int) (float64, error)
		GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error)
	}

	// WalletRepository defines the methods for wallet operations
	WalletQueryRepositoryHandler interface {
		// Wallet CRUD operations
		GetAllWallets(ctx context.Context) ([]entity.Wallet, error)
		GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error)
		// Balance & Transactions
		GetBalance(ctx context.Context, walletID int) (float64, error)
		GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error)
	}
)
