-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE,

    order_id    BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    amount      DECIMAL(10,2) NOT NULL CHECK (amount > 0),
    currency    TEXT NOT NULL DEFAULT 'USD',
    method      TEXT NOT NULL CHECK (method IN ('credit_card', 'debit_card', 'paypal', 'gift_card', 'store_credit')),
    status      TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'succeeded', 'failed', 'refunded')),
    transaction_id TEXT UNIQUE,
    gateway_response JSONB
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_user_id ON payments(user_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payments;
-- +goose StatementEnd
