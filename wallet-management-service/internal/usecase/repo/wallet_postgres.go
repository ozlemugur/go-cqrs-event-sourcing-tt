package repo

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/postgres"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/entity"
)

type WalletRepo struct {
	*postgres.Postgres
}

// NewWalletRepo creates a new instance of WalletRepo
func NewWalletRepo(pg *postgres.Postgres) *WalletRepo {
	return &WalletRepo{pg}
}

// GetWalletByID retrieves a wallet by its ID (returns active wallets only)
func (r *WalletRepo) GetWalletByID(ctx context.Context, id int) (*entity.WalletResponse, error) {
	sql, args, err := r.Builder.
		Select("id, address, network, status, created_at").
		From("wallets").
		Where("id = ? AND status = ?", id, "active").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletByID - Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	var wallet entity.WalletResponse
	err = row.Scan(&wallet.ID, &wallet.Address, &wallet.Network, &wallet.Status, &wallet.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletByID - Scan: %w", err)
	}

	return &wallet, nil
}

// CreateWallet inserts a new wallet into the database
func (r *WalletRepo) CreateWallet(ctx context.Context, wallet entity.WalletRequest) error {
	sql, args, err := r.Builder.
		Insert("wallets").
		Columns("address", "network", "status").
		Values(wallet.Address, wallet.Network, "active").
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletRepo - CreateWallet - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - CreateWallet - Exec: %w", err)
	}

	return nil
}

// UpdateWallet updates an existing wallet in the database
func (r *WalletRepo) UpdateWallet(ctx context.Context, id int, wallet entity.WalletRequest) error {
	sql, args, err := r.Builder.
		Update("wallets").
		Set("address", wallet.Address).
		Set("network", wallet.Network).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletRepo - UpdateWallet - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - UpdateWallet - Exec: %w", err)
	}

	return nil
}

// DeleteWallet performs a soft delete by updating the wallet status to "deleted"
func (r *WalletRepo) DeleteWallet(ctx context.Context, id int) error {
	sql, args, err := r.Builder.
		Update("wallets").
		Set("status", "deleted").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletRepo - DeleteWallet - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - DeleteWallet - Exec: %w", err)
	}

	return nil
}
