package entity

import "time"

// WalletRequest represents the request payload when creating or updating a wallet
type WalletRequest struct {
	Address string `json:"address" binding:"required"` // Wallet address (required)
	Network string `json:"network" binding:"required"` // Network type (required)
}

// WalletResponse represents the response payload when retrieving a wallet
type WalletResponse struct {
	ID        int       `json:"id" db:"id"`           // Wallet ID
	Address   string    `json:"address" db:"address"` // Wallet address
	Network   string    `json:"network" db:"network"` // Network type
	Status    string    `json:"status" db:"status"`   // Status field (optional)
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
