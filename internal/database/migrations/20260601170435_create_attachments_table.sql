-- +goose Up
CREATE TABLE attachment(
    attachment_id INTEGER PRIMARY KEY,
    message_id INTEGER NOT NULL,
    file_url TEXT NOT NULL,
    file_type TEXT NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages(message_id)
);


-- +goose Down
DROP TABLE attachment;
