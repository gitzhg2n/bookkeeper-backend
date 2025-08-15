-- +migrate Up
CREATE TABLE IF NOT EXISTS bills (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    amount_cents BIGINT NOT NULL,
    due_day INT NOT NULL,
    next_due TIMESTAMP NOT NULL,
    recurring BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +migrate Down
DROP TABLE IF EXISTS bills;
