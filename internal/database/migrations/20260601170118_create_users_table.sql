-- +goose Up
CREATE TABLE IF NOT EXISTS users(
   user_id INTEGER PRIMARY KEY,
   username TEXT,
   password_hash TEXT NOT NULL,
   phone_number TEXT NOT NULL,
   CONSTRAINT users_uc_phone UNIQUE (phone_number)
);

-- +goose Down
DROP TABLE users;
