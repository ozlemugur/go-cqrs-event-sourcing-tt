package usecase

import (
	"context"
	"fmt"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/wallet-management-service/internal/entity"
)

// WalletUseCase - İş mantığı için WalletHandler implementasyonu
type WalletUseCase struct {
	repo WalletRepositoryHandler
	l    logger.Interface
}

// NewWalletUseCase - Yeni bir WalletUseCase instance'ı oluşturur.
func NewWalletUseCase(r WalletRepositoryHandler, l logger.Interface) *WalletUseCase {
	return &WalletUseCase{
		repo: r,
		l:    l,
	}
}

// GetAllWallets - Veritabanından tüm cüzdanları getirir
func (uc *WalletUseCase) GetAllWallets(ctx context.Context) ([]entity.Wallet, error) {
	wallets, err := uc.repo.GetAllWallets(ctx)
	if err != nil {
		return nil, fmt.Errorf("WalletUseCase - GetAllWallets - uc.repo.GetAllWallets: %w", err)
	}
	return wallets, nil
}

// GetWalletByID - Belirtilen ID'ye sahip cüzdanı getirir
func (uc *WalletUseCase) GetWalletByID(ctx context.Context, id int) (*entity.Wallet, error) {
	wallet, err := uc.repo.GetWalletByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("WalletUseCase - GetWalletByID - uc.repo.GetWalletByID: %w", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("WalletUseCase - GetWalletByID: wallet not found")
	}
	return wallet, nil
}

// CreateWallet - Yeni bir cüzdan oluşturur
func (uc *WalletUseCase) CreateWallet(ctx context.Context, wallet entity.Wallet) error {
	// Address ve Network boş olmamalı
	if wallet.Address == "" || wallet.Network == "" {
		return fmt.Errorf("WalletUseCase - CreateWallet: address or network is empty")
	}

	err := uc.repo.CreateWallet(ctx, wallet)
	if err != nil {
		return fmt.Errorf("WalletUseCase - CreateWallet - uc.repo.CreateWallet: %w", err)
	}
	return nil
}

// UpdateWallet - Belirtilen ID'ye sahip bir cüzdanı günceller
func (uc *WalletUseCase) UpdateWallet(ctx context.Context, id int, wallet entity.Wallet) error {
	// Address ve Network boş olmamalı
	if wallet.Address == "" || wallet.Network == "" {
		return fmt.Errorf("WalletUseCase - UpdateWallet: address or network is empty")
	}

	err := uc.repo.UpdateWallet(ctx, id, wallet)
	if err != nil {
		return fmt.Errorf("WalletUseCase - UpdateWallet - uc.repo.UpdateWallet: %w", err)
	}
	return nil
}

// DeleteWallet - Belirtilen ID'ye sahip bir cüzdanı siler
func (uc *WalletUseCase) DeleteWallet(ctx context.Context, id int) error {
	err := uc.repo.DeleteWallet(ctx, id)
	if err != nil {
		return fmt.Errorf("WalletUseCase - DeleteWallet - uc.repo.DeleteWallet: %w", err)
	}
	return nil
}
