// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/entity"
)

type (

	// WalletQueryHandler defines methods for querying wallet and asset data
	// WalletQueryUseCaseHandler defines the interface for wallet asset operations
	WalletQueryUseCaseHandler interface {
		// Retrieves all assets and their balances for a specific wallet ID
		GetAllAssets(ctx context.Context, walletID int) ([]entity.WalletAsset, error)

		// Retrieves the balance of a specific asset in a wallet
		GetAssetBalance(ctx context.Context, walletID int, assetName string) (*entity.WalletAsset, error)

		// Retrieves the transaction history for a specific wallet and asset
		GetTransactionHistory(ctx context.Context, walletID int, assetName string) ([]entity.Transaction, error)
	}

	// WalletQueryRepositoryHandler defines the methods for querying wallet data.
	WalletQueryRepositoryHandler interface {
		// GetAssetsByWalletID retrieves all assets and their amounts for a specific wallet ID.
		GetAssetsByWalletID(ctx context.Context, walletID int) ([]entity.WalletAsset, error)

		// GetWalletAsset retrieves the amount of a specific asset for a wallet.
		GetWalletAsset(ctx context.Context, walletID int, assetName string) (float64, error)

		// UpdateWalletAsset updates the amount of a specific asset for a wallet.
		UpdateWalletAsset(ctx context.Context, walletID int, assetName string, amount float64) error

		// InsertOrUpdateWalletAsset inserts or updates a wallet asset entry.
		InsertOrUpdateWalletAsset(ctx context.Context, walletID int, assetName string, amount float64) error

		// GetTransactionHistory retrieves transaction history for a wallet.
		GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error)
	}
)
