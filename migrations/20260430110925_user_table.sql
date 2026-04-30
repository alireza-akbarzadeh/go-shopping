-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id                BIGSERIAL PRIMARY KEY,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMP WITH TIME ZONE,

    email             TEXT NOT NULL UNIQUE,
    email_verified_at TIMESTAMP WITH TIME ZONE,

    phone             TEXT,
    first_name        TEXT NOT NULL,
    last_name         TEXT NOT NULL,

    password_hash     TEXT NOT NULL,
    role              TEXT NOT NULL DEFAULT 'user',
    is_active         BOOLEAN NOT NULL DEFAULT TRUE,

    last_login_at     TIMESTAMP WITH TIME ZONE
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);

-- Comments (optional, for clarity)
COMMENT ON COLUMN users.id IS 'Auto‑incrementing primary key';
COMMENT ON COLUMN users.created_at IS 'Creation timestamp';
COMMENT ON COLUMN users.updated_at IS 'Last update timestamp';
COMMENT ON COLUMN users.deleted_at IS 'Soft‑delete timestamp (null = active)';
COMMENT ON COLUMN users.email IS 'Unique login email';
COMMENT ON COLUMN users.email_verified_at IS 'When email was verified (null = unverified)';
COMMENT ON COLUMN users.phone IS 'Phone number in E.164 format';
COMMENT ON COLUMN users.first_name IS 'Given name';
COMMENT ON COLUMN users.last_name IS 'Family name';
COMMENT ON COLUMN users.password_hash IS 'bcrypt (or argon2) hash – never returned to client';
COMMENT ON COLUMN users.role IS 'Role: user, admin, moderator';
COMMENT ON COLUMN users.is_active IS 'Soft disable flag (soft‑delete is separate)';
COMMENT ON COLUMN users.last_login_at IS 'Timestamp of last successful login';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
