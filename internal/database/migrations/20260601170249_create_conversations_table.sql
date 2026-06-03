-- +goose Up
CREATE TABLE conversations(
    conversation_id INTEGER PRIMARY KEY,
    type_of TEXT NOT NULL,
    title TEXT NOT NULL,
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES messages(message_id),
    FOREIGN KEY (parent_id) REFERENCES conversations(conversation_id)
);
-- +goose Down
DROP TABLE conversations;