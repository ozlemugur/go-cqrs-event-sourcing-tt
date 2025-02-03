package usecase

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/entity"
)

// WalletUseCase implements the business logic for wallet operations.
type WalletUseCase struct {
	repo WalletRepositoryHandler
	log  logger.Interface
}

// NewWalletUseCase creates a new instance of WalletUseCase.
func NewWalletUseCase(r WalletRepositoryHandler, log logger.Interface) *WalletUseCase {
	return &WalletUseCase{
		repo: r,
		log:  log,
	}
}

// GetWalletByID retrieves a wallet by its ID.
func (uc *WalletUseCase) GetWalletByID(ctx context.Context, id int) (*entity.WalletResponse, error) {
	wallet, err := uc.repo.GetWalletByID(ctx, id)
	if err != nil {
		uc.log.Error(err, "WalletUseCase - GetWalletByID - repository error")
		return nil, fmt.Errorf("WalletUseCase - GetWalletByID: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("WalletUseCase - GetWalletByID: wallet not found")
	}
	return wallet, nil
}

// CreateWallet adds a new wallet to the database.
func (uc *WalletUseCase) CreateWallet(ctx context.Context, wallet entity.WalletRequest) error {
	// Validate required fields
	if wallet.Address == "" || wallet.Network == "" {
		return fmt.Errorf("WalletUseCase - CreateWallet: address or network is empty")
	}

	if err := uc.repo.CreateWallet(ctx, wallet); err != nil {
		uc.log.Error(err, "WalletUseCase - CreateWallet - repository error")
		return fmt.Errorf("WalletUseCase - CreateWallet: %w", err)
	}

	uc.log.Info("Wallet created successfully", "address", wallet.Address)
	return nil
}

// UpdateWallet updates an existing wallet by its ID.
func (uc *WalletUseCase) UpdateWallet(ctx context.Context, id int, wallet entity.WalletRequest) error {
	// Validate required fields
	if wallet.Address == "" || wallet.Network == "" {
		return fmt.Errorf("WalletUseCase - UpdateWallet: address or network is empty")
	}

	if err := uc.repo.UpdateWallet(ctx, id, wallet); err != nil {
		uc.log.Error(err, "WalletUseCase - UpdateWallet - repository error")
		return fmt.Errorf("WalletUseCase - UpdateWallet: %w", err)
	}

	uc.log.Info("Wallet updated successfully", "wallet_id", id)
	return nil
}

// DeleteWallet performs a soft delete of a wallet by its ID.
func (uc *WalletUseCase) DeleteWallet(ctx context.Context, id int) error {
	if err := uc.repo.DeleteWallet(ctx, id); err != nil {
		uc.log.Error(err, "WalletUseCase - DeleteWallet - repository error")
		return fmt.Errorf("WalletUseCase - DeleteWallet: %w", err)
	}

	uc.log.Info("Wallet deleted successfully", "wallet_id", id)
	return nil
}
