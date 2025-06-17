CREATE TABLE IF NOT EXISTS outbox (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL,
    published_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT now()
);
