-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shipments (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE,

    order_id    BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    user_id     BIGINT NOT NULL REFERENCES users(id),
    carrier     TEXT NOT NULL,
    tracking_number TEXT,
    status      TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','shipped','delivered','returned')),
    shipped_at  TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    estimated_delivery TIMESTAMP WITH TIME ZONE,

    address_line1 TEXT NOT NULL,
    address_line2 TEXT,
    city         TEXT NOT NULL,
    state        TEXT,
    postal_code  TEXT NOT NULL,
    country      TEXT NOT NULL
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_user_id ON shipments(user_id);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
CREATE INDEX idx_shipments_status ON shipments(status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shipments;
-- +goose StatementEnd
