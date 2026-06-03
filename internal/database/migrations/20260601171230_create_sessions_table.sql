-- +goose Up
CREATE TABLE sessions(
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX idx_sessions_expiry ON sessions(expiry);
-- +goose Down
DROP TABLE sessions;