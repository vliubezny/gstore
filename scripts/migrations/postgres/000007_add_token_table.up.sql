CREATE TABLE IF NOT EXISTS token (
    id UUID PRIMARY KEY,
    user_id INTEGER REFERENCES store_user (id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL
);