CREATE EXTENSION pgcrypto;

CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    account_id UUID NOT NULL,
    encrypted_pan BYTEA NOT NULL,
    hashed_pan BYTEA NOT NULL,
    expiry_month INT NOT NULL CHECK (expiry_month >= 1 AND expiry_month <= 12),
    expiry_year INT NOT NULL,
    -- encrypted_cvv VARCHAR(255) NOT NULL,
    cardholder_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- CREATE UNIQUE INDEX IF NOT EXISTS uidx_card_number ON cards(number);
