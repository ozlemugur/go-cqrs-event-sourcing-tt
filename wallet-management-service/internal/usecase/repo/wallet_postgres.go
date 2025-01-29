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

// GetAllWallets retrieves all wallets from the database
func (r *WalletRepo) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
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
func (r *WalletRepo) GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error) {
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

// CreateWallet inserts a new wallet into the database
func (r *WalletRepo) CreateWallet(ctx context.Context, wallet entity.Wallet) error {
	sql, args, err := r.Builder.
		Insert("wallets").
		Columns("address", "network").
		Values(wallet.Address, wallet.Network).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletRepo - CreateWallet - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - CreateWallet - r.Pool.Exec: %w", err)
	}

	return nil
}

// UpdateWallet updates an existing wallet in the database
func (r *WalletRepo) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) error {
	sql, args, err := r.Builder.
		Update("wallets").
		Set("address", wallet.Address).
		Set("network", wallet.Network).
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletRepo - UpdateWallet - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - UpdateWallet - r.Pool.Exec: %w", err)
	}

	return nil
}

// DeleteWallet deletes a wallet by its ID
func (r *WalletRepo) DeleteWallet(ctx context.Context, id int) error {
	sql, args, err := r.Builder.
		Delete("wallets").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("WalletRepo - DeleteWallet - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - DeleteWallet - r.Pool.Exec: %w", err)
	}

	return nil
}
