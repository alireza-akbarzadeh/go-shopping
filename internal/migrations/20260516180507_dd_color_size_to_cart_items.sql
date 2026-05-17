-- +goose Up
-- +goose StatementBegin
ALTER TABLE cart_items 
ADD COLUMN color VARCHAR(50) NOT NULL DEFAULT '',
ADD COLUMN size VARCHAR(50) NOT NULL DEFAULT '';

-- Optional: composite index for faster variant lookups
CREATE INDEX idx_cart_items_product_variant ON cart_items(product_id, color, size);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cart_items_product_variant;
ALTER TABLE cart_items DROP COLUMN color, DROP COLUMN size;
-- +goose StatementEnd