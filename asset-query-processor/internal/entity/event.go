package entity

import "time"

// Event struct represents a wallet transaction event from Kafka
type Event struct {
	EventID   string  `json:"event_id"`
	WalletID  int     `json:"wallet_id"`
	Type      string  `json:"type"` // "withdraw", "deposit", "transfer"
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
}

// WalletEvent represents a single event in the event store.{"event_id":"c0b0b54a-6c49-4848-9b6a-1413564056c1","wallet_id":1,"asset_name":"BTC","type":"deposit","amount":100,"timestamp":1738492914}%
type WalletEvent struct {
	EventID   string  `json:"event_id" db:"event_id"`           // Unique identifier for the event
	WalletID  int     `json:"wallet_id" db:"wallet_id"`         // Wallet associated with this event
	AssetName string  `json:"asset_name" db:"asset_name"`       // asset name
	Type      string  `json:"type" db:"type"`                   // Event type: "withdraw", "deposit", "transfer"
	Amount    float64 `json:"amount" db:"amount"`               // Transaction amount
	Timestamp int64   `json:"timestamp" db:"timestamp"`         // Unix timestamp when the event was created
	Metadata  string  `json:"metadata,omitempty" db:"metadata"` // Optional JSON metadata (for extensibility)
}

// Transaction represents a financial operation on a wallet.
type Transaction struct {
	ID        int       `json:"id" db:"transaction_id"`
	WalletID  int       `json:"wallet_id" db:"wallet_id"`
	Type      string    `json:"type" db:"type"` // Possible values: "withdraw", "deposit", "transfer"
	Amount    float64   `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

/*
// NewWalletEvent creates a new WalletEvent instance.
func NewWalletEvent(walletID int, eventType string, amount float64, metadata string) WalletEvent {
	return WalletEvent{
		EventID:   GenerateUUID(), // Unique ID for event tracking
		WalletID:  walletID,
		Type:      eventType,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
		Metadata:  metadata,
	}
} */
