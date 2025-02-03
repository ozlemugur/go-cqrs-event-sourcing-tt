CREATE TABLE IF NOT EXISTS  wallets (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    network VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',    -- Status of the wallet (e.g., active, deleted)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(address, network)
);

