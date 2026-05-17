-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS slug TEXT UNIQUE,
    ADD COLUMN IF NOT EXISTS compare_at_price DECIMAL(10,2) CHECK (compare_at_price >= 0),
    ADD COLUMN IF NOT EXISTS cost DECIMAL(10,2) CHECK (cost >= 0),
    ADD COLUMN IF NOT EXISTS barcode TEXT,
    ADD COLUMN IF NOT EXISTS low_stock_threshold INT DEFAULT 5,
    ADD COLUMN IF NOT EXISTS weight DECIMAL(8,2) CHECK (weight >= 0),
    ADD COLUMN IF NOT EXISTS is_digital BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS meta_title TEXT,
    ADD COLUMN IF NOT EXISTS meta_description TEXT;

-- Update existing rows: set slug from id as temporary default
UPDATE products SET slug = 'product-' || id WHERE slug IS NULL;

ALTER TABLE products ALTER COLUMN slug SET NOT NULL;
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_is_digital ON products(is_digital);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_products_slug;
DROP INDEX IF EXISTS idx_products_is_digital;
ALTER TABLE products
    DROP COLUMN IF EXISTS slug,
    DROP COLUMN IF EXISTS compare_at_price,
    DROP COLUMN IF EXISTS cost,
    DROP COLUMN IF EXISTS barcode,
    DROP COLUMN IF EXISTS low_stock_threshold,
    DROP COLUMN IF EXISTS weight,
    DROP COLUMN IF EXISTS is_digital,
    DROP COLUMN IF EXISTS meta_title,
    DROP COLUMN IF EXISTS meta_description;
-- +goose StatementEnd
