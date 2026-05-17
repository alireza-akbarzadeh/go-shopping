-- +goose Up
-- +goose StatementBegin

-- Table: menu_groups
CREATE TABLE menu_groups (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Table: menu_items (self-referential for nested items)
CREATE TABLE menu_items (
    id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES menu_groups(id) ON DELETE CASCADE,
    parent_id INT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    label VARCHAR(255) NOT NULL,
    href VARCHAR(500) NULL,
    icon VARCHAR(100) NOT NULL,
    permission VARCHAR(100) NULL,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_menu_items_group_id ON menu_items(group_id);
CREATE INDEX idx_menu_items_parent_id ON menu_items(parent_id);
CREATE INDEX idx_menu_items_permission ON menu_items(permission);
CREATE INDEX idx_menu_groups_display_order ON menu_groups(display_order);
CREATE INDEX idx_menu_items_display_order ON menu_items(display_order);

-- +goose StatementEnd