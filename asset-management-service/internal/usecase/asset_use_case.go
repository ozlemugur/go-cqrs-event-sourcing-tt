package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-management-service/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

// AssetUseCase handles asset transactions using event sourcing.
type AssetUseCase struct {
	commandQueue CommandProducerHandler // Event producer (Kafka)
	log          logger.Interface
}

// NewAssetUseCase creates a new asset use case.
func NewAssetUseCase(commandQueue CommandProducerHandler, l logger.Interface) *AssetUseCase {
	return &AssetUseCase{
		commandQueue: commandQueue,
		log:          l,
	}
}

// Withdraw funds from a wallet (Publishes event)
func (uc *AssetUseCase) Withdraw(ctx context.Context, walletID int, assetName string, amount float64) error {
	command := entity.Command{
		CommandID: uuid.New().String(),
		WalletID:  walletID,
		Type:      "withdraw",
		AssetName: assetName,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	if err := uc.commandQueue.PublishCommand(ctx, command); err != nil {
		return fmt.Errorf("Withdraw - PublishEvent: %w", err)
	}

	uc.log.Info("Withdraw event published", "WalletID", walletID, "AssetName", assetName, "Amount", amount)
	return nil
}

// Deposit funds into a wallet (Publishes event)
func (uc *AssetUseCase) Deposit(ctx context.Context, walletID int, assetName string, amount float64) error {
	//wallet check , ya header üzerinden teyitli geldiğini var sayabiliriz ya da httpcall ve ya readonly bir check yapabiliriz
	command := entity.Command{
		CommandID: uuid.New().String(),
		WalletID:  walletID,
		Type:      "deposit",
		AssetName: assetName,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}

	if err := uc.commandQueue.PublishCommand(ctx, command); err != nil {
		return fmt.Errorf("Deposit - PublishEvent: %w", err)
	}

	uc.log.Info("Deposit event published", "WalletID", walletID, "AssetName", assetName, "Amount", amount)
	return nil
}

// Transfer funds between wallets (Publishes a TransferCommand)
func (uc *AssetUseCase) Transfer(ctx context.Context, fromWalletID, toWalletID int, assetName string, amount float64, executeTime int64) error {
	if executeTime == 0 || executeTime <= time.Now().Unix() {
		executeTime = time.Now().Unix()
	}
	status := "scheduled"
	if executeTime <= time.Now().Unix() {
		status = "executed" // Eğer executeTime geçmişte bir zamansa, işlemi hemen gerçekleştir
	}
	transferCommand := entity.TransferCommand{
		CommandID:   uuid.New().String(), // Benzersiz komut ID'si
		FromWallet:  fromWalletID,
		ToWallet:    toWalletID,
		AssetName:   assetName,
		Amount:      amount,
		Type:        "transfer",
		ExecuteTime: executeTime, // Komutun yürütüleceği zaman zamanı
		Status:      status,      // Başlangıç durumu
		CreatedAt:   time.Now(),  // Oluşturulma zamanı
	}

	// Komut Kafka'ya gönderiliyor
	if err := uc.commandQueue.PublishCommand(ctx, transferCommand); err != nil {
		return fmt.Errorf("Transfer - PublishCommand: %w", err)
	}

	uc.log.Info("Transfer command published",
		"FromWalletID", fromWalletID,
		"ToWalletID", toWalletID,
		"AssetName", assetName,
		"Amount", amount,
	)
	return nil
}
