CREATE TYPE order_status AS ENUM
(
  'pending',
  'confirmed',
  'paid',
  'shipped',
  'delivered',
  'cancelled'
);

CREATE TABLE orders
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status     order_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);