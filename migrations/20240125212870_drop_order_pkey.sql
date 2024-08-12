-- +goose Up
ALTER TABLE orders DROP CONSTRAINT orders_pkey;

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders ADD PRIMARY KEY (id);
-- +goose StatementEnd
