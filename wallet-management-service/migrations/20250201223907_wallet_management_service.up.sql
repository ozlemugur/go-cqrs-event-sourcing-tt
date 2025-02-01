CREATE TABLE IF NOT EXISTS  wallets (
    id SERIAL PRIMARY KEY,
    address VARCHAR(255) NOT NULL,
    network VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(address, network)
);

CREATE TABLE IF NOT EXISTS  wallet_assets (
    id SERIAL PRIMARY KEY,
    wallet_id INT NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    amount DECIMAL(18, 8) NOT NULL DEFAULT 0,
    FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE,
    UNIQUE(wallet_id, asset_name)
);