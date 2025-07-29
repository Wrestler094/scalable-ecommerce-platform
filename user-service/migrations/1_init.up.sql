CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_email ON users (email);

ALTER TABLE users
    ADD CONSTRAINT check_user_role CHECK (role IN ('user', 'admin'));
