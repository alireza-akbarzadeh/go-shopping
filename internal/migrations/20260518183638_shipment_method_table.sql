-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shipping_methods (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true
);

-- Create indexes for better query performance
CREATE INDEX idx_shipping_methods_is_active ON shipping_methods(is_active);
CREATE INDEX idx_shipping_methods_name ON shipping_methods(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_shipping_methods_is_active;
DROP INDEX IF EXISTS idx_shipping_methods_name;
DROP TABLE IF EXISTS shipping_methods;
-- +goose StatementEnd