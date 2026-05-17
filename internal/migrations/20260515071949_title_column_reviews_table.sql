-- +goose Up
-- +goose StatementBegin
ALTER TABLE reviews
    ADD COLUMN IF NOT EXISTS title TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE reviews
    DROP COLUMN IF EXISTS title;
-- +goose StatementEnd