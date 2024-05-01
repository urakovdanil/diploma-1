-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD CONSTRAINT unique_user_order UNIQUE (user_id, number);
ALTER TABLE orders
    ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;

CREATE OR REPLACE FUNCTION update_orders_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER orders_update_trigger
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION update_orders_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER orders_update_trigger ON orders;
DROP FUNCTION IF EXISTS update_orders_updated_at();

ALTER TABLE orders
    DROP COLUMN created_at,
    DROP COLUMN updated_at;
ALTER TABLE orders DROP CONSTRAINT unique_user_order;

-- +goose StatementEnd
