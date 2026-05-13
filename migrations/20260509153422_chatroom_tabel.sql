-- +goose Up
CREATE TABLE chat_rooms (
    room_id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INTEGER NOT NULL,
    is_private BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_chat_rooms_created_by ON chat_rooms(created_by);
CREATE INDEX idx_chat_rooms_is_private ON chat_rooms(is_private);

-- +goose Down
DROP INDEX IF EXISTS idx_chat_rooms_is_private;
DROP INDEX IF EXISTS idx_chat_rooms_created_by;
DROP TABLE IF EXISTS chat_rooms;
