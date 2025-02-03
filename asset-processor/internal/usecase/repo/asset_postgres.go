package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/postgres"
)

type AssetRepo struct {
	*postgres.Postgres
}

// NewAssetRepo - Creates a new repository instance
func NewAssetRepo(pg *postgres.Postgres) *AssetRepo {
	return &AssetRepo{pg}
}

// GetBalance - Retrieves the balance of a wallet
func (r *AssetRepo) GetBalance(ctx context.Context, walletID int) (float64, error) {
	var balance float64
	sql, args, err := r.Builder.
		Select("balance").
		From("wallets").
		Where("id = ?", walletID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("AssetRepo - GetBalance - Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("AssetRepo - GetBalance - QueryRow: %w", err)
	}

	return balance, nil
}

// UpdateBalance - Updates the balance of a wallet
func (r *AssetRepo) UpdateBalance(ctx context.Context, walletID int, newBalance float64) error {
	sql, args, err := r.Builder.
		Update("wallets").
		Set("balance", newBalance).
		Where("id = ?", walletID).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetRepo - UpdateBalance - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AssetRepo - UpdateBalance - Exec: %w", err)
	}

	return nil
}

// InsertTransaction - Stores a new transaction in the database
func (r *AssetRepo) InsertTransaction(ctx context.Context, transaction entity.Transaction) error {
	sql, args, err := r.Builder.
		Insert("transactions").
		Columns("wallet_id", "type", "amount", "status", "created_at").
		Values(transaction.WalletID, transaction.Type, transaction.Amount, transaction.Status, transaction.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetRepo - InsertTransaction - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AssetRepo - InsertTransaction - Exec: %w", err)
	}

	return nil
}

// InsertScheduledTransaction - Stores a scheduled transaction
func (r *AssetRepo) InsertScheduledTransaction(ctx context.Context, transaction entity.ScheduledTransaction) error {
	sql, args, err := r.Builder.
		Insert("scheduled_transactions").
		Columns("from_wallet", "to_wallet", "amount", "execute_time", "status", "created_at").
		Values(transaction.FromWallet, transaction.ToWallet, transaction.Amount, transaction.ExecuteTime, transaction.Status, transaction.CreatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetRepo - InsertScheduledTransaction - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AssetRepo - InsertScheduledTransaction - Exec: %w", err)
	}

	return nil
}

// GetScheduledTransactions - Retrieves all scheduled transactions due for execution
func (r *AssetRepo) GetScheduledTransactions(ctx context.Context, executeTime time.Time) ([]entity.ScheduledTransaction, error) {
	sql, args, err := r.Builder.
		Select("id, from_wallet, to_wallet, amount, execute_time, status, created_at").
		From("scheduled_transactions").
		Where("execute_time <= ? AND status = ?", executeTime, "scheduled").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("AssetRepo - GetScheduledTransactions - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("AssetRepo - GetScheduledTransactions - Query: %w", err)
	}
	defer rows.Close()

	transactions := make([]entity.ScheduledTransaction, 0)

	for rows.Next() {
		var txn entity.ScheduledTransaction
		err = rows.Scan(&txn.ID, &txn.FromWallet, &txn.ToWallet, &txn.Amount, &txn.ExecuteTime, &txn.Status, &txn.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("AssetRepo - GetScheduledTransactions - Scan: %w", err)
		}

		transactions = append(transactions, txn)
	}

	return transactions, nil
}

// MarkTransactionAsProcessed - Marks a scheduled transaction as processed
func (r *AssetRepo) MarkTransactionAsProcessed(ctx context.Context, transactionID int) error {
	sql, args, err := r.Builder.
		Update("scheduled_transactions").
		Set("status", "executed").
		Where("id = ?", transactionID).
		ToSql()
	if err != nil {
		return fmt.Errorf("AssetRepo - MarkTransactionAsProcessed - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AssetRepo - MarkTransactionAsProcessed - Exec: %w", err)
	}

	return nil
}
