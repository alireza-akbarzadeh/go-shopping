-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS rating DECIMAL(3,2) DEFAULT 0.0,
    ADD COLUMN IF NOT EXISTS reviews_count INT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_new BOOLEAN DEFAULT FALSE;

-- Automatically set is_new based on creation date (optional)
UPDATE products SET is_new = (created_at > NOW() - INTERVAL '30 days');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS reviews_count,
    DROP COLUMN IF EXISTS is_new;
-- +goose StatementEnd