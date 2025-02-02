package repo

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/postgres"
)

type WalletQueryRepo struct {
	*postgres.Postgres
}

// NewWalletQueryRepo creates a new instance of WalletQueryRepo
func NewWalletQueryRepo(pg *postgres.Postgres) *WalletQueryRepo {
	return &WalletQueryRepo{pg}
}

// GetAssetsByWalletID retrieves all assets and their amounts for a specific wallet ID
func (r *WalletQueryRepo) GetAssetsByWalletID(ctx context.Context, walletID int) ([]entity.WalletAsset, error) {
	sql, _, err := r.Builder.
		Select("asset_name, amount").
		From("wallet_assets").
		Where("wallet_id = ?", walletID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletQueryRepo - GetAssetsByWalletID - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryRepo - GetAssetsByWalletID - Query: %w", err)
	}
	defer rows.Close()

	assets := make([]entity.WalletAsset, 0)
	for rows.Next() {
		var asset entity.WalletAsset
		err = rows.Scan(&asset.AssetName, &asset.Amount)
		if err != nil {
			return nil, fmt.Errorf("WalletQueryRepo - GetAssetsByWalletID - Scan: %w", err)
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// GetWalletAsset retrieves the amount of a specific asset for a wallet
func (r *WalletQueryRepo) GetWalletAsset(ctx context.Context, walletID int, assetName string) (float64, error) {
	var amount float64

	sql, _, err := r.Builder.
		Select("amount").
		From("wallet_assets").
		Where("wallet_id = ? AND asset_name = ?", walletID, assetName).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("WalletQueryRepo - GetWalletAsset - Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, walletID, assetName).Scan(&amount)
	if err != nil {
		return 0, fmt.Errorf("WalletQueryRepo - GetWalletAsset - QueryRow: %w", err)
	}

	return amount, nil
}

// UpdateWalletAsset updates the amount of a specific asset for a wallet
func (r *WalletQueryRepo) UpdateWalletAsset(ctx context.Context, walletID int, assetName string, amount float64) error {
	sql, _, err := r.Builder.
		Update("wallet_assets").
		Set("amount", amount).
		Where("wallet_id = ? AND asset_name = ?", walletID, assetName).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletQueryRepo - UpdateWalletAsset - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, walletID, assetName)
	if err != nil {
		return fmt.Errorf("WalletQueryRepo - UpdateWalletAsset - Exec: %w", err)
	}

	return nil
}

// InsertOrUpdateWalletAsset inserts or updates a wallet asset entry
func (r *WalletQueryRepo) InsertOrUpdateWalletAsset(ctx context.Context, walletID int, assetName string, amount float64) error {
	sql, _, err := r.Builder.
		Insert("wallet_assets").
		Columns("wallet_id", "asset_name", "amount").
		Values(walletID, assetName, amount).
		Suffix("ON CONFLICT (wallet_id, asset_name) DO UPDATE SET amount = EXCLUDED.amount").
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletQueryRepo - InsertOrUpdateWalletAsset - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, walletID, assetName, amount)
	if err != nil {
		return fmt.Errorf("WalletQueryRepo - InsertOrUpdateWalletAsset - Exec: %w", err)
	}

	return nil
}

// GetTransactionHistory retrieves transaction history for a wallet
func (r *WalletQueryRepo) GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error) {
	sql, _, err := r.Builder.
		Select("transaction_id, wallet_id, asset_name, type, amount, created_at").
		From("wallet_transactions").
		Where("wallet_id = ?", walletID).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletQueryRepo - GetTransactionHistory - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletQueryRepo - GetTransactionHistory - Query: %w", err)
	}
	defer rows.Close()

	transactions := make([]entity.Transaction, 0)
	for rows.Next() {
		var txn entity.Transaction
		err = rows.Scan(&txn.ID, &txn.WalletID, &txn.AssetName, &txn.Type, &txn.Amount, &txn.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("WalletQueryRepo - GetTransactionHistory - Scan: %w", err)
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}
