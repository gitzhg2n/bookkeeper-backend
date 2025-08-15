-- +migrate Up
CREATE TABLE IF NOT EXISTS investment_alerts (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    asset_symbol VARCHAR(32) NOT NULL,
    alert_type VARCHAR(32) NOT NULL,
    direction VARCHAR(8) NOT NULL,
    threshold DOUBLE PRECISION NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS investment_alerts;
