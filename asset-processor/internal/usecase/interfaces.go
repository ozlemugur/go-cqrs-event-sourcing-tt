package usecase

import (
	"context"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/entity"
)

type (
	/* Asset Management UseCase Interface */
	AssetUseCaseHandler interface {
		Withdraw(ctx context.Context, walletID int, assetName string, amount float64) error // Withdraw funds from a wallet
		Deposit(ctx context.Context, walletID int, assetName string, amount float64) error  // Deposit funds into a wallet
	}

	// EventJournal defines the contract for publishing events.
	EventJournal interface {
		PublishEvent(ctx context.Context, event entity.WalletEvent) error
	}

	/* Command Handler  UseCase Interface */
	CommandHandler interface {
		//ProcessCommand(ctx context.Context, command entity.Command) error
		MsgfessageHandler(key, value []byte) error
	}

	RetryEventProducer interface {
		PublishRetryEvent(ctx context.Context, event []byte) error
	}

	DLQEventProducer interface {
		PublishDLQEvent(ctx context.Context, event []byte) error
	}
	// Command Queue defines the contract for publishing events.
	CommandProducerHandler interface {
		PublishCommand(ctx context.Context, command interface{}) error
	}
)
