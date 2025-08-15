-- +migrate Up
CREATE TABLE IF NOT EXISTS budgets (
    id SERIAL PRIMARY KEY,
    household_id BIGINT NOT NULL,
    month VARCHAR(7) NOT NULL,
    category_id BIGINT NOT NULL,
    planned_cents BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL,
    user_id BIGINT,
    amount_cents BIGINT NOT NULL,
    currency VARCHAR(8) NOT NULL,
    category_id BIGINT,
    memo TEXT,
    occurred_at TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS budgets;
