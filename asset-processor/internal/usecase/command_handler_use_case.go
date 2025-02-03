package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/asset-processor/internal/entity"
	"github.com/ozlemugur/go-cqrs-event-sourcing-tt/pkg/logger"
)

type commandHandler struct {
	log          logger.Interface
	assetUseCase AssetUseCaseHandler
	commandQueue CommandProducerHandler
	handlers     map[string]func(ctx context.Context, data []byte) error
}

// NewCommandHandler initializes the command handler and maps command types to their corresponding functions.
func NewCommandHandler(commandQueue CommandProducerHandler, assetUseCase AssetUseCaseHandler, l logger.Interface) *commandHandler {
	h := &commandHandler{
		log:          l,
		assetUseCase: assetUseCase,
		commandQueue: commandQueue,
	}

	// Initialize command handlers map
	h.handlers = map[string]func(ctx context.Context, data []byte) error{
		"withdraw": h.handleWithdrawCommand,
		"deposit":  h.handleDepositCommand,
		"transfer": h.handleTransferCommand,
	}

	return h
}

// MsgfessageHandler is the main entry point for handling messages.
func (h *commandHandler) MsgfessageHandler(key, value []byte) error {
	h.log.Info("Received message", "key", string(key))

	if len(value) == 0 {
		return fmt.Errorf("empty message value")
	}

	ctx := context.Background()

	// Decode Base64 message
	decodedValue, err := base64.StdEncoding.DecodeString(strings.Trim(string(value), "\""))
	if err != nil {
		h.log.Error(err, "Base64 decode error")
		return err
	}

	// Extract the base command to determine its type
	var baseCommand entity.BaseCommand
	if err := json.Unmarshal(decodedValue, &baseCommand); err != nil {
		h.log.Error(err, "Failed to unmarshal BaseCommand")
		return err
	}

	// Find and execute the appropriate handler
	if handler, exists := h.handlers[baseCommand.Type]; exists {
		return handler(ctx, decodedValue)
	}

	return fmt.Errorf("unknown command type: %s", baseCommand.Type)
}

// **Command Handlers**

func (h *commandHandler) handleWithdrawCommand(ctx context.Context, data []byte) error {
	var command entity.Command
	if err := json.Unmarshal(data, &command); err != nil {
		h.log.Error(err, "Failed to unmarshal withdraw command")
		return err
	}

	h.log.Info("Processing withdraw command",
		"WalletID", command.WalletID,
		"AssetName", command.AssetName,
		"Amount", command.Amount,
	)

	if err := h.assetUseCase.Withdraw(ctx, command.WalletID, command.AssetName, command.Amount); err != nil {
		h.log.Error(err, "Withdraw failed")
		return fmt.Errorf("handleWithdrawCommand: %w", err)
	}

	h.log.Info("Withdraw command processed successfully")
	return nil
}

func (h *commandHandler) handleDepositCommand(ctx context.Context, data []byte) error {
	var command entity.Command
	if err := json.Unmarshal(data, &command); err != nil {
		h.log.Error(err, "Failed to unmarshal deposit command")
		return err
	}

	h.log.Info("Processing deposit command",
		"WalletID", command.WalletID,
		"AssetName", command.AssetName,
		"Amount", command.Amount,
	)

	if err := h.assetUseCase.Deposit(ctx, command.WalletID, command.AssetName, command.Amount); err != nil {
		h.log.Error(err, "Deposit failed")
		return fmt.Errorf("handleDepositCommand: %w", err)
	}

	h.log.Info("Deposit command processed successfully")
	return nil
}

func (h *commandHandler) handleTransferCommand(ctx context.Context, data []byte) error {
	var command entity.TransferCommand
	if err := json.Unmarshal(data, &command); err != nil {
		h.log.Error(err, "Failed to unmarshal transfer command")
		return err
	}

	h.log.Info("Processing transfer command",
		"FromWallet", command.FromWallet,
		"ToWallet", command.ToWallet,
		"AssetName", command.AssetName,
		"Amount", command.Amount,
		"ExecuteTime", command.ExecuteTime,
	)

	currentTime := time.Now().Unix()

	// Re-schedule if ExecuteTime is too far in the future
	if command.ExecuteTime > currentTime+600 {
		h.log.Info("Execute time is more than 10 minutes in the future. Re-scheduling command",
			"CommandID", command.CommandID,
		)
		if err := h.commandQueue.PublishCommand(ctx, command); err != nil {
			h.log.Error(err, "Failed to re-schedule transfer command")
			return fmt.Errorf("handleTransferCommand - Retry PublishCommand: %w", err)
		}
		return nil
	}

	if err := h.assetUseCase.Withdraw(ctx, command.FromWallet, command.AssetName, command.Amount); err != nil {
		h.log.Error(err, "Withdraw failed")
		return fmt.Errorf("handleTransferCommand - Withdraw: %w", err)
	}

	if err := h.assetUseCase.Deposit(ctx, command.ToWallet, command.AssetName, command.Amount); err != nil {
		h.log.Error(err, "Deposit failed")
		return fmt.Errorf("handleTransferCommand - Deposit: %w", err)
	}

	h.log.Info("Transfer command processed successfully")
	return nil
}
