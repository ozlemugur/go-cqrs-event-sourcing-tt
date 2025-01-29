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

// NewWalletRepo creates a new instance of WalletRepo
func NewWalletQueryRepo(pg *postgres.Postgres) *WalletQueryRepo {
	return &WalletQueryRepo{pg}
}

// GetAllWallets retrieves all wallets from the database
func (r *WalletQueryRepo) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	sql, _, err := r.Builder.
		Select("id, address, network, created_at").
		From("wallets").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetAllWallets - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetAllWallets - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	wallets := make([]entity.Wallet, 0)
	for rows.Next() {
		wallet := entity.Wallet{}
		err = rows.Scan(&wallet.ID, &wallet.Address, &wallet.Network, &wallet.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("WalletRepo - GetAllWallets - rows.Scan: %w", err)
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// GetWalletByID retrieves a wallet by its ID
func (r *WalletQueryRepo) GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error) {
	sql, _, err := r.Builder.
		Select("id, address, network, created_at").
		From("wallets").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletByID - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, id)

	wallet := &entity.Wallet{}
	err = row.Scan(&wallet.ID, &wallet.Address, &wallet.Network, &wallet.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletByID - row.Scan: %w", err)
	}

	return wallet, nil
}

// GetBalance retrieves the balance of a wallet from the read model
func (r *WalletQueryRepo) GetBalance(ctx context.Context, walletID int) (float64, error) {
	var balance float64

	sql, _, err := r.Builder.
		Select("balance").
		From("wallet_balances").
		Where("wallet_id = ?", walletID).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("WalletRepo - GetBalance - Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("WalletRepo - GetBalance - QueryRow: %w", err)
	}

	return balance, nil
}

// GetTransactionHistory retrieves transaction history for a wallet
func (r *WalletQueryRepo) GetTransactionHistory(ctx context.Context, walletID int) ([]entity.Transaction, error) {
	sql, _, err := r.Builder.
		Select("transaction_id, wallet_id, type, amount, created_at").
		From("wallet_transactions").
		Where("wallet_id = ?", walletID).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetTransactionHistory - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, walletID)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetTransactionHistory - Query: %w", err)
	}
	defer rows.Close()

	transactions := make([]entity.Transaction, 0)
	for rows.Next() {
		var txn entity.Transaction
		err = rows.Scan(&txn.ID, &txn.WalletID, &txn.Type, &txn.Amount, &txn.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("WalletRepo - GetTransactionHistory - Scan: %w", err)
		}
		transactions = append(transactions, txn)
	}

	return transactions, nil
}
