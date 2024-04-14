-- +goose Up
ALTER TABLE queue
   ADD COLUMN   info jsonb DEFAULT '{}'::jsonb;

-- +goose Down
-- +goose StatementBegin
ALTER TABLE queue
    DROP COLUMN   info;
-- +goose StatementEnd
