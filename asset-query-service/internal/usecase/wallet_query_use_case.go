package usecase

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

// WalletQueryUseCase implements business logic for wallet asset operations
type WalletQueryUseCase struct {
	repo WalletQueryRepositoryHandler
	log  logger.Interface
}

// NewWalletQueryUseCase creates a new instance of WalletQueryUseCase
func NewWalletQueryUseCase(r WalletQueryRepositoryHandler, log logger.Interface) *WalletQueryUseCase {
	return &WalletQueryUseCase{repo: r, log: log}
}

// GetAllAssets retrieves all assets for a given wallet ID
func (uc *WalletQueryUseCase) GetAllAssets(ctx context.Context, walletID int) ([]entity.WalletAsset, error) {
	assets, err := uc.repo.GetAssetsByWalletID(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetAllAssets - uc.repo.GetAssetsByWalletID: %w", err)
	}
	return assets, nil
}

// GetAssetBalance retrieves the balance of a specific asset in a wallet
func (uc *WalletQueryUseCase) GetAssetBalance(ctx context.Context, walletID int, assetName string) (*entity.WalletAsset, error) {
	amount, err := uc.repo.GetWalletAsset(ctx, walletID, assetName)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetAssetBalance - uc.repo.GetWalletAsset: %w", err)
	}
	return &entity.WalletAsset{
		WalletID:  walletID,
		AssetName: assetName,
		Amount:    amount,
	}, nil
}

// GetTransactionHistory retrieves the transaction history for a given wallet and asset
func (uc *WalletQueryUseCase) GetTransactionHistory(ctx context.Context, walletID int, assetName string) ([]entity.Transaction, error) {
	transactions, err := uc.repo.GetTransactionHistory(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryUseCase - GetTransactionHistory - uc.repo.GetTransactionHistory: %w", err)
	}
	return filterTransactionsByAsset(transactions, assetName), nil
}

// Helper function to filter transactions by asset name
func filterTransactionsByAsset(transactions []entity.Transaction, assetName string) []entity.Transaction {
	filtered := make([]entity.Transaction, 0)
	for _, txn := range transactions {
		if txn.AssetName == assetName {
			filtered = append(filtered, txn)
		}
	}
	return filtered
}
