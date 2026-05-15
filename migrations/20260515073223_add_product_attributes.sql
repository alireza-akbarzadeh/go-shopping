-- +goose Up
-- +goose StatementBegin
ALTER TABLE products
    ADD COLUMN IF NOT EXISTS colors JSONB DEFAULT '[]',
    ADD COLUMN IF NOT EXISTS sizes JSONB DEFAULT '[]';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE products
    DROP COLUMN IF EXISTS colors,
    DROP COLUMN IF EXISTS sizes;
-- +goose StatementEnd