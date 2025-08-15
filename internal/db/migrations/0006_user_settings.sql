-- +migrate Up
CREATE TABLE IF NOT EXISTS user_settings (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    large_transaction_threshold BIGINT NOT NULL DEFAULT 25000
);

-- +migrate Down
DROP TABLE IF EXISTS user_settings;
