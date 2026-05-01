-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE,

    name        TEXT NOT NULL,
    description TEXT,
    price       DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    stock       INT NOT NULL DEFAULT 0 CHECK (stock >= 0),
    sku         TEXT NOT NULL UNIQUE,
    category_id BIGINT,
    images      TEXT[],
    status      TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'archived')),
    created_by  BIGINT REFERENCES users(id) ON DELETE SET NULL,
    updated_by  BIGINT REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_products_deleted_at ON products(deleted_at);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_category_id ON products(category_id);
CREATE INDEX idx_products_status ON products(status);

COMMENT ON COLUMN products.id IS 'Product unique identifier';
COMMENT ON COLUMN products.name IS 'Product name';
COMMENT ON COLUMN products.description IS 'Detailed product description';
COMMENT ON COLUMN products.price IS 'Current price in USD (decimal)';
COMMENT ON COLUMN products.stock IS 'Available quantity';
COMMENT ON COLUMN products.sku IS 'Stock Keeping Unit – unique';
COMMENT ON COLUMN products.category_id IS 'Optional reference to categories table (to be added later)';
COMMENT ON COLUMN products.images IS 'Array of image URLs';
COMMENT ON COLUMN products.status IS 'active, inactive, archived';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
