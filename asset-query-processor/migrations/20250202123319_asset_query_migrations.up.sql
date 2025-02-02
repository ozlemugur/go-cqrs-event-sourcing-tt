CREATE TABLE IF NOT EXISTS  wallet_assets (
    wallet_id INT NOT NULL,
    asset_name VARCHAR(50) NOT NULL,
    amount FLOAT NOT NULL DEFAULT 0,
    PRIMARY KEY (wallet_id, asset_name)
);

