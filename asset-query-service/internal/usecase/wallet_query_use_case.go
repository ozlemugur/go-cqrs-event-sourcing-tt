package usecase

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

// WalletQueryUseCase - Implements business logic for wallet operations
type WalletQueryUseCase struct {
	repo WalletQueryRepositoryHandler
	l    logger.Interface
}

// NewWalletQueryUseCase - Creates a new instance of WalletQueryUseCase
func NewWalletQueryUseCase(r WalletQueryRepositoryHandler, l logger.Interface) *WalletQueryUseCase {
	return &WalletQueryUseCase{
		repo: r,
		l:    l,
	}
}

// GetAllWallets - Retrieves all wallets from the database
func (uc *WalletQueryUseCase) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	wallets, err := uc.repo.GetAllWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetAllWallets - uc.repo.GetAllWallets: %w", err)
	}
	return wallets, nil
}

// GetWalletByID - Retrieves a wallet by its ID
func (uc *WalletQueryUseCase) GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error) {
	wallet, err := uc.repo.GetWalletByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetWalletByID - uc.repo.GetWalletByID: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetWalletByID: wallet not found")
	}
	return wallet, nil
}

// GetBalance - Retrieves the current balance of a wallet by its ID
func (uc *WalletQueryUseCase) GetBalance(ctx context.Context, walletID int) (float64, error) {
	balance, err := uc.repo.GetBalance(ctx, walletID)
	if err != nil {
		return 0, fmt.Errorf("WalletQueryUseCase - GetBalance - uc.repo.GetBalance: %w", err)
	}
	return balance, nil
}

// GetTransactionHistory - Retrieves transaction history for a specific wallet
func (uc *WalletQueryUseCase) GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error) {
	transactions, err := uc.repo.GetTransactionHistory(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetTransactionHistory - uc.repo.GetTransactionHistory: %w", err)
	}
	return transactions, nil
}
