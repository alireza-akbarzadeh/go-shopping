-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMP WITH TIME ZONE,

    name        TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    description TEXT,
    parent_id   BIGINT REFERENCES categories(id) ON DELETE CASCADE,
    level       INT NOT NULL DEFAULT 0,
    path        TEXT,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE INDEX idx_categories_deleted_at ON categories(deleted_at);
CREATE INDEX idx_categories_slug ON categories(slug);
CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_level ON categories(level);
CREATE INDEX idx_categories_path ON categories(path);

COMMENT ON COLUMN categories.parent_id IS 'Self-reference to parent category (NULL = root)';
COMMENT ON COLUMN categories.level IS 'Nesting depth (0 = root)';
COMMENT ON COLUMN categories.path IS 'Materialized path for descendant queries';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
