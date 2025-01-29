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

// GetBalance - Retrieves the current balance of a wallet.
func (r *AssetQueryRepo) GetBalance(ctx context.Context, walletID int) (float64, error) {
	var balance float64

	sql, _, err := r.Builder.
		Select("balance").
		From("wallet_balances").
		Where("wallet_id = ?", walletID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("AssetQueryRepo - GetBalance - Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("AssetQueryRepo - GetBalance - QueryRow: %w", err)
	}

	return balance, nil
}

// UpdateBalance - Updates the balance of a wallet.
func (r *AssetQueryRepo) UpdateBalance(ctx context.Context, walletID int, newBalance float64) error {
	sql, _, err := r.Builder.
		Update("wallet_balances").
		Set("balance", newBalance).
		Where("wallet_id = ?", walletID).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - UpdateBalance - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, walletID)
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - UpdateBalance - Exec: %w", err)
	}

	return nil
}

// InsertTransaction - Stores a new transaction in the database.
func (r *AssetQueryRepo) InsertTransaction(ctx context.Context, txn entity.Transaction) error {
	sql, _, err := r.Builder.
		Insert("wallet_transactions").
		Columns("wallet_id", "type", "amount", "created_at").
		Values(txn.WalletID, txn.Type, txn.Amount, txn.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - InsertTransaction - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("AssetQueryRepo - InsertTransaction - Exec: %w", err)
	}

	return nil
}

// GetTransactionHistory - Retrieves all transactions for a wallet.
func (r *AssetQueryRepo) GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error) {
	sql, _, err := r.Builder.
		Select("transaction_id, wallet_id, type, amount, created_at").
		From("wallet_transactions").
		Where("wallet_id = ?", walletID).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AssetQueryRepo - GetTransactionHistory - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, walletID)
	if err != nil {
		return nil, fmt.Errorf("AssetQueryRepo - GetTransactionHistory - Query: %w", err)
	}
	defer rows.Close()

	transactions := make([]entity.Transaction, 0)
	for rows.Next() {
		var txn entity.Transaction
		err = rows.Scan(&txn.ID, &txn.WalletID, &txn.Type, &txn.Amount, &txn.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("AssetQueryRepo - GetTransactionHistory - Scan: %w", err)
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}
