-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS coupons (
    id                   BIGSERIAL PRIMARY KEY,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at           TIMESTAMP WITH TIME ZONE,

    code                 TEXT NOT NULL UNIQUE,
    description          TEXT,
    discount_type        TEXT NOT NULL CHECK (discount_type IN ('percentage', 'fixed')),
    discount_value       DECIMAL(10,2) NOT NULL CHECK (discount_value > 0),
    minimum_order_amount DECIMAL(10,2) DEFAULT 0 CHECK (minimum_order_amount >= 0),
    max_discount_amount  DECIMAL(10,2), -- only for percentage discounts
    usage_limit          INT DEFAULT 1 CHECK (usage_limit > 0),
    used_count           INT NOT NULL DEFAULT 0 CHECK (used_count >= 0),
    start_date           TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date             TIMESTAMP WITH TIME ZONE NOT NULL,
    is_active            BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_coupons_code ON coupons(code);
CREATE INDEX idx_coupons_active ON coupons(is_active);
CREATE INDEX idx_coupons_date_range ON coupons(start_date, end_date);

CREATE TABLE IF NOT EXISTS coupon_usages (
    id               BIGSERIAL PRIMARY KEY,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    coupon_id        BIGINT NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
    user_id          BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    order_id         BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    discount_amount  DECIMAL(10,2) NOT NULL CHECK (discount_amount > 0)
);

CREATE INDEX idx_coupon_usages_coupon ON coupon_usages(coupon_id);
CREATE INDEX idx_coupon_usages_user ON coupon_usages(user_id);
CREATE INDEX idx_coupon_usages_order ON coupon_usages(order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS coupon_usages;
DROP TABLE IF EXISTS coupons;
-- +goose StatementEnd