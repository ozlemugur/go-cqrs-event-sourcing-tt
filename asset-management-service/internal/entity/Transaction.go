package entity

import (
	"time"
)

// TransactionRequest represents a request for withdraw/deposit operations
type TransactionRequest struct {
	WalletID  int     `json:"wallet_id" binding:"required"`  // Wallet performing the transaction
	AssetName string  `json:"asset_name" binding:"required"` // Asset being transacted (e.g., BTC, ETH)
	Amount    float64 `json:"amount" binding:"required"`     // Transaction amount
}

// TransferRequest represents a request for transferring funds between wallets
type TransferRequest struct {
	FromWalletID int     `json:"from_wallet_id" binding:"required"` // Sender wallet
	ToWalletID   int     `json:"to_wallet_id" binding:"required"`   // Receiver wallet
	AssetName    string  `json:"asset_name" binding:"required"`     // Asset being transferred
	Amount       float64 `json:"amount" binding:"required"`         // Transfer amount
}

// ScheduledTransactionRequest represents a request for scheduling future transactions
type ScheduledTransactionRequest struct {
	FromWalletID int     `json:"from_wallet_id" binding:"required"` // Sender wallet
	ToWalletID   int     `json:"to_wallet_id" binding:"required"`   // Receiver wallet
	AssetName    string  `json:"asset_name" binding:"required"`     // Asset to transfer
	Amount       float64 `json:"amount" binding:"required"`         // Amount to transfer
	ExecuteTime  int64   `json:"execute_time" binding:"required"`   // Execution time (Unix timestamp)
}

// WalletEvent represents an event in the event journal
type WalletEvent struct {
	EventID   string  `json:"event_id"`   // Unique event identifier
	WalletID  int     `json:"wallet_id"`  // Associated wallet
	AssetName string  `json:"asset_name"` // Asset being transacted
	Type      string  `json:"type"`       // "withdraw", "deposit", "transfer"
	Amount    float64 `json:"amount"`     // Transaction amount
	Timestamp int64   `json:"timestamp"`  // Event time (Unix)
}

// ScheduledTransactionEvent represents an event for scheduled transactions
type ScheduledTransactionEvent struct {
	EventID     int     `json:"event_id"`     // Unique scheduled transaction ID
	WalletID    int     `json:"wallet_id"`    // Wallet related to the transaction
	AssetName   string  `json:"asset_name"`   // Asset for the scheduled transaction
	Type        string  `json:"type"`         // "scheduled_withdraw", "scheduled_deposit"
	Amount      float64 `json:"amount"`       // Amount for the scheduled transaction
	ExecuteTime int64   `json:"execute_time"` // Execution timestamp (Unix format)
}

// Transaction represents a completed transaction record
type Transaction struct {
	ID        int       `json:"id" db:"id"`
	WalletID  int       `json:"wallet_id" db:"wallet_id"`
	AssetName string    `json:"asset_name" db:"asset_name"`
	Type      string    `json:"type" db:"type"` // "withdraw" | "deposit" | "transfer"
	Amount    float64   `json:"amount" db:"amount"`
	Status    string    `json:"status" db:"status"` // "pending", "completed", "failed"
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ScheduledTransaction represents a future-dated transaction
type ScheduledTransaction struct {
	ID          int       `json:"id" db:"id"`
	FromWallet  int       `json:"from_wallet" db:"from_wallet"`
	ToWallet    int       `json:"to_wallet" db:"to_wallet"`
	AssetName   string    `json:"asset_name" db:"asset_name"`
	Amount      float64   `json:"amount" db:"amount"`
	ExecuteTime time.Time `json:"execute_time" db:"execute_time"` // When this should be executed
	Status      string    `json:"status" db:"status"`             // "scheduled", "executed"
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
