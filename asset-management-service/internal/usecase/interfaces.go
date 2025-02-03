package usecase

import (
	"context"
)

type (
	/* Asset Management UseCase Interface */
	AssetHandler interface {
		Withdraw(ctx context.Context, walletID int, assetName string, amount float64) error                                        // Withdraw funds from a wallet
		Deposit(ctx context.Context, walletID int, assetName string, amount float64) error                                         // Deposit funds into a wallet
		Transfer(ctx context.Context, fromWalletID int, toWalletID int, assetName string, amount float64, executeTime int64) error // Transfer funds between wallets
	}

	// EventJournal defines the contract for publishing events.
	CommandProducerHandler interface {
		PublishCommand(ctx context.Context, command interface{}) error
	}
)
