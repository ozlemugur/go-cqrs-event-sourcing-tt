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

// GetAllWallets retrieves all wallets along with their assets from the database
func (r *WalletRepo) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	sql, _, err := r.Builder.
		Select("w.id, w.address, w.network, w.created_at, a.asset_name, a.amount").
		From("wallets AS w").
		LeftJoin("wallet_assets AS a ON w.id = a.wallet_id").
		OrderBy("w.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetAllWallets - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetAllWallets - Query: %w", err)
	}
	defer rows.Close()

	wallets := make(map[int]*entity.Wallet)
	for rows.Next() {
		var wallet entity.Wallet
		var asset entity.Asset

		err := rows.Scan(&wallet.ID, &wallet.Address, &wallet.Network, &wallet.CreatedAt, &asset.Name, &asset.Amount)
		if err != nil {
			return nil, fmt.Errorf("WalletRepo - GetAllWallets - Scan: %w", err)
		}

		if existingWallet, found := wallets[wallet.ID]; found {
			existingWallet.Assets = append(existingWallet.Assets, asset)
		} else {
			wallet.Assets = []entity.Asset{asset}
			wallets[wallet.ID] = &wallet
		}
	}

	result := make([]entity.Wallet, 0, len(wallets))
	for _, wallet := range wallets {
		result = append(result, *wallet)
	}

	return result, nil
}

// GetWalletByID retrieves a wallet by its ID along with its assets
func (r *WalletRepo) GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error) {
	sql, _, err := r.Builder.
		Select("w.id, w.address, w.network, w.created_at, a.asset_name, a.amount").
		From("wallets AS w").
		LeftJoin("wallet_assets AS a ON w.id = a.wallet_id").
		Where("w.id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletByID - Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("WalletRepo - GetWalletByID - Query: %w", err)
	}
	defer rows.Close()

	var wallet *entity.Wallet
	for rows.Next() {
		if wallet == nil {
			wallet = &entity.Wallet{}
		}
		var asset entity.Asset

		err := rows.Scan(&wallet.ID, &wallet.Address, &wallet.Network, &wallet.CreatedAt, &asset.Name, &asset.Amount)
		if err != nil {
			return nil, fmt.Errorf("WalletRepo - GetWalletByID - Scan: %w", err)
		}

		wallet.Assets = append(wallet.Assets, asset)
	}

	if wallet == nil {
		return nil, nil // Wallet not found
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
		return fmt.Errorf("WalletRepo - CreateWallet - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - CreateWallet - Exec: %w", err)
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
		return fmt.Errorf("WalletRepo - UpdateWallet - Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("WalletRepo - UpdateWallet - Exec: %w", err)
	}

	return nil
}

// DeleteWallet deletes a wallet and its associated assets by its ID
func (r *WalletRepo) DeleteWallet(ctx context.Context, id int) error {
	// Start a transaction to ensure atomic delete for both wallet and assets
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("WalletRepo - DeleteWallet - Begin: %w", err)
	}

	// Delete assets associated with the wallet
	sql, args, err := r.Builder.
		Delete("wallet_assets").
		Where("wallet_id = ?", id).
		ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("WalletRepo - DeleteWallet - DeleteAssets - Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("WalletRepo - DeleteWallet - DeleteAssets - Exec: %w", err)
	}

	// Delete the wallet
	sql, args, err = r.Builder.
		Delete("wallets").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("WalletRepo - DeleteWallet - DeleteWallet - Builder: %w", err)
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("WalletRepo - DeleteWallet - DeleteWallet - Exec: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("WalletRepo - DeleteWallet - Commit: %w", err)
	}

	return nil
}
