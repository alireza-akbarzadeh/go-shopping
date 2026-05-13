-- +goose Up
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    message_type VARCHAR(50) DEFAULT 'text',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES chat_rooms(room_id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_room_id ON messages(room_id);
CREATE INDEX idx_messages_user_id ON messages(user_id);
CREATE INDEX idx_messages_created_at ON messages(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_user_id;
DROP INDEX IF EXISTS idx_messages_room_id;
DROP TABLE IF EXISTS messages;
