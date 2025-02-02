package entity

import "time"

// Transaction represents a wallet transaction (e.g., withdraw, deposit, transfer).
type Transaction struct {
	ID             int       `json:"id" db:"transaction_id"`                           // Unique transaction ID
	WalletID       int       `json:"wallet_id" db:"wallet_id"`                         // Wallet associated with the transaction
	TargetWalletID *int      `json:"target_wallet_id,omitempty" db:"target_wallet_id"` // For transfer transactions
	Type           string    `json:"type" db:"type"`                                   // "withdraw", "deposit", or "transfer"
	AssetName      string    `json:"asset_name" db:"asset_name"`
	Amount         float64   `json:"amount" db:"amount"`         // Transaction amount
	CreatedAt      time.Time `json:"created_at" db:"created_at"` // Timestamp when the transaction occurred
}
