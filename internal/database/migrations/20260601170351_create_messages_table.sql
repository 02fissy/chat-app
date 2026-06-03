-- +goose Up
CREATE TABLE messages(
    message_id INTEGER PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    conversation_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    parent_id INTEGER NOT NULL,
    FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id)
);

-- +goose Down
DROP TABLE messages;