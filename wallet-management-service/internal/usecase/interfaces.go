// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/entity"
)

type (
	/* Wallet UseCase Interface */
	WalletHandler interface {
		// CRUD operations
		GetWalletByID(ctx context.Context, id int) (*entity.WalletResponse, error)   // Fetch a wallet by ID
		CreateWallet(ctx context.Context, wallet entity.WalletRequest) error         // Create a new wallet
		UpdateWallet(ctx context.Context, id int, wallet entity.WalletRequest) error // Update an existing wallet
		DeleteWallet(ctx context.Context, id int) error                              // Delete a wallet by ID
	}

	/* Wallet Repository Interface */
	WalletRepositoryHandler interface {
		// Database layer methods for wallets
		GetWalletByID(ctx context.Context, id int) (*entity.WalletResponse, error)   // Fetch a wallet by ID
		CreateWallet(ctx context.Context, wallet entity.WalletRequest) error         // Insert a wallet into the database
		UpdateWallet(ctx context.Context, id int, wallet entity.WalletRequest) error // Update a wallet record in the database
		DeleteWallet(ctx context.Context, id int) error                              // Delete a wallet record
	}
)
