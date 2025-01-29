package entity

import "time"

type Wallet struct {
	ID        int       `json:"-" db:"id"`
	Address   string    `json:"address"` // Wallet address
	Network   string    `json:"network"` // Wallet network
	CreatedAt time.Time `json:"-" db:"created_at"`
}
