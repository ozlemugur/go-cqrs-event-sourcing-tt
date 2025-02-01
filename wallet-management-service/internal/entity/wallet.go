package entity

import "time"

type Wallet struct {
	ID        int       `json:"id"`
	Address   string    `json:"address"`
	Network   string    `json:"network"`
	Assets    []Asset   `json:"assets,omitempty"`
	CreatedAt time.Time `json:"-" db:"created_at"`
}

type Asset struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}
