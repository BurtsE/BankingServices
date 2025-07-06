
CREATE TYPE currency_type AS ENUM ('rub', 'usd');

CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(256) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
    balance NUMERIC(18,4) NOT NULL DEFAULT 0,
    currency currency_type NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);