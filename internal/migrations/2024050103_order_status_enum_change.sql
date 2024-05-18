-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status_v2 AS ENUM ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED', 'NEW');
ALTER TABLE orders ALTER COLUMN status DROP DEFAULT;
ALTER TABLE orders ALTER COLUMN status TYPE order_status_v2 USING status::text::order_status_v2;
ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'NEW'::order_status_v2;
DROP TYPE order_status;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('REGISTERED', 'INVALID', 'PROCESSING', 'PROCESSED');
ALTER TABLE orders ALTER COLUMN status DROP DEFAULT;
UPDATE orders SET status = 'PROCESSING' WHERE status = 'NEW';
ALTER TABLE orders ALTER COLUMN status TYPE order_status USING status::text::order_status;
ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'PROCESSING'::order_status_v2;
DROP TYPE order_status_v2;

-- +goose StatementEnd
