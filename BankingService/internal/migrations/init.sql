
CREATE TYPE currency_type AS ENUM ('rub', 'usd');

CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(256) NOT NULL,
    balance NUMERIC(18,4) NOT NULL DEFAULT 0,
    currency currency_type NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount NUMERIC(18,4) NOT NULL,
    currency currency_type NOT NULL,
    type VARCHAR(32) NOT NULL,     -- deposit, withdraw, transfer, payment, etc.
    status VARCHAR(32) NOT NULL,   -- success, pending, failed, etc.
    description TEXT,
    related_entity_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);