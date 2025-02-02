package repo

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-query-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/postgres"
)

// AssetQueryRepo provides methods for querying balances & transactions.
type AssetQueryRepo struct {
	*postgres.Postgres
}

// NewAssetQueryRepo - Creates a new repository instance.
func NewAssetQueryRepo(pg *postgres.Postgres) *AssetQueryRepo {
	return &AssetQueryRepo{pg}
}

// GetBalance - Retrieves the current balance of a wallet's asset.
func (r *AssetQueryRepo) GetBalance(ctx context.Context, walletID int, assetName string) (float64, error) {
	var balance float64

	sql, _, err := r.Builder.
		Select("amount").
		From("wallet_assets").
		Where("wallet_id = ? AND asset_name = ?", walletID, assetName).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("AssetQueryRepo - GetBalance - Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, walletID, assetName).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("AssetQueryRepo - GetBalance - QueryRow: %w", err)
	}

	return balance, nil
}

// UpdateBalance - Updates the balance of a wallet's asset.
func (r *AssetQueryRepo) UpdateBalance(ctx context.Context, walletID int, assetName string, amount float64) error {
	sql := `
	INSERT INTO wallet_assets (wallet_id, asset_name, amount)
	VALUES ($1, $2, $3)
	ON CONFLICT (wallet_id, asset_name) DO UPDATE
	SET amount = wallet_assets.amount + $3
	`

	_, err := r.Pool.Exec(ctx, sql, walletID, assetName, amount)
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - UpdateBalance - Exec: %w", err)
	}

	return nil
}

// InsertTransaction - Stores a new transaction in the database.
func (r *AssetQueryRepo) InsertTransaction(ctx context.Context, txn entity.Transaction) error {
	/*sql, _, err := r.Builder.
		Insert("wallet_transactions").
		Columns("wallet_id", "asset_name", "type", "amount", "created_at").
		Values(txn.WalletID, txn.AssetName, txn.Type, txn.Amount, txn.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - InsertTransaction - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - InsertTransaction - Exec: %w", err)
	} */

	return nil
}

// GetTransactionHistory - Retrieves all transactions for a wallet and asset.
func (r *AssetQueryRepo) GetTransactionHistory(ctx context.Context, walletID int, assetName string) ([]entity.Transaction, error) {
	/*sql, _, err := r.Builder.
		Select("transaction_id, wallet_id, asset_name, type, amount, created_at").
		From("wallet_transactions").
		Where("wallet_id = ? AND asset_name = ?", walletID, assetName).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AssetQueryRepo - GetTransactionHistory - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, walletID, assetName)
	if err != nil {
		return nil, fmt.Errorf("AssetQueryRepo - GetTransactionHistory - Query: %w", err)
	}
	defer rows.Close()

	transactions := make([]entity.Transaction, 0)
	for rows.Next() {
		var txn entity.Transaction
		err = rows.Scan(&txn.ID, &txn.WalletID, &txn.AssetName, &txn.Type, &txn.Amount, &txn.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("AssetQueryRepo - GetTransactionHistory - Scan: %w", err)
		}
		transactions = append(transactions, txn)
	} */

	return nil, nil
}
