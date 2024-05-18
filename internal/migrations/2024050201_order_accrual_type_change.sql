-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE float USING accrual::float;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ALTER COLUMN accrual TYPE integer USING round(accrual)::integer;

-- +goose StatementEnd
