-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    number TEXT NOT NULL,
    status order_status NOT NULL DEFAULT 'PROCESSING'::order_status,
    accrual INTEGER NOT NULL DEFAULT 0,
    user_id INTEGER REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE order_status;

-- +goose StatementEnd
