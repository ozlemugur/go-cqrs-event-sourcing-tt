package entity

import "time"

/*type Wallet struct {
	ID        int       `json:"-" db:"id"`
	Address   string    `json:"address"` // Wallet address
	Network   string    `json:"network"` // Wallet network
	CreatedAt time.Time `json:"-" db:"created_at"`
} */

// WalletAsset represents the balance of a specific asset within a wallet.
type WalletAsset struct {
	WalletID  int       `json:"wallet_id" db:"wallet_id"`             // The ID of the wallet
	AssetName string    `json:"asset_name" db:"asset_name"`           // The name of the asset (e.g., BTC, ETH)
	Amount    float64   `json:"amount" db:"amount"`                   // The amount of the asset
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"` // Timestamp of the last update
}
