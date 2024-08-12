-- +goose Up
ALTER TABLE orders
    ADD COLUMN   info jsonb DEFAULT '{}'::jsonb;

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN   info;
-- +goose StatementEnd
