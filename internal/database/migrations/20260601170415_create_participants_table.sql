-- +goose Up
CREATE TABLE participants(
    participant_id INTEGER PRIMARY KEY,
    role TEXT NOT NULL,
    conversation_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    last_read_message_id INTEGER,
    FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (last_read_message_id) REFERENCES messages(message_id)
);

-- +goose Down
DROP TABLE participants;
