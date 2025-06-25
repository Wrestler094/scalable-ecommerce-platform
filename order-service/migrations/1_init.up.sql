CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY, -- внутренний числовой идентификатор
    uuid UUID NOT NULL DEFAULT gen_random_uuid(), -- внешний идентификатор
    user_id BIGINT NOT NULL,
    status TEXT NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (uuid)
);

CREATE TABLE IF NOT EXISTS order_items (
    order_id BIGINT REFERENCES orders(id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    PRIMARY KEY (order_id, product_id)
);