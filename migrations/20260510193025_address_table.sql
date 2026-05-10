-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS addresses (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMP WITH TIME ZONE,

    user_id          BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    address_type     TEXT NOT NULL CHECK (address_type IN ('shipping', 'billing', 'both')),
    is_default       BOOLEAN NOT NULL DEFAULT FALSE,
    recipient_name   TEXT NOT NULL,
    phone            TEXT NOT NULL,
    address_line1    TEXT NOT NULL,
    address_line2    TEXT,
    city             TEXT NOT NULL,
    state            TEXT,
    postal_code      TEXT NOT NULL,
    country          TEXT NOT NULL,
    instructions     TEXT
);

CREATE INDEX idx_addresses_user_id ON addresses(user_id);
CREATE INDEX idx_addresses_type ON addresses(address_type);
CREATE INDEX idx_addresses_default ON addresses(is_default);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS addresses;
-- +goose StatementEnd