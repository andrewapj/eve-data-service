CREATE TABLE IF NOT EXISTS expiry
(
    kind VARCHAR(16) PRIMARY KEY,
    expires TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
