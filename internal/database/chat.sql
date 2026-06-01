CREATE TABLE users(
   user_id INTEGER PRIMARY KEY,
   username TEXT NOT NULL,
   password_hash TEXT NOT NULL,
   phone_number TEXT NOT NULL,
   CONSTRAINT users_uc_phone UNIQUE (phone_number);
);


CREATE TABLE contact(
    contact_id INTEGER PRIMARY KEY,
    user1_id INTEGER NOT NULL,
    user2_id INTEGER NOT NULL,
    FOREIGN KEY (user1_id) REFERENCES users(user_id),
    FOREIGN KEY (user2_id) REFERENCES users(user_id),
    UNIQUE(user1_id, user2_id)
);
CREATE TABLE conversations(
    conversation_id INTEGER PRIMARY KEY,
    type_of TEXT NOT NULL,
    title TEXT NOT NULL,
    parent_id INTEGER,
    FOREIGN KEY (parent_id) REFERENCES messages(message_id),
    FOREIGN KEY (parent_id) REFERENCES conversations(conversation_id)
);
CREATE TABLE messages(
    message_id INTEGER PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    conversation_id INTEGER NOT NULL,
    FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id),
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    parent_id INTEGER NOT NULL
);
CREATE TABLE participants(
    participant_id INTEGER PRIMARY KEY,
    role TEXT NOT NULL,
    conversation_id INTEGER NOT NULL,
    FOREIGN KEY (conversation_id) REFERENCES conversations(conversation_id),
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    last_read_message_id INTEGER,
    FOREIGN KEY (last_read_message_id) REFERENCES messages(message_id)
);
CREATE TABLE attachment(
    attachment_id INTEGER PRIMARY KEY,
    message_id INTEGER NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages(message_id),
    file_url TEXT NOT NULL,
    file_type TEXT NOT NULL
);

CREATE TABLE sessions(
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);
CREATE INDEX idx_sessions_expiry ON sessions(expiry);