-- +goose Up
ALTER TABLE savepoints
    ALTER COLUMN value DROP DEFAULT,
    ALTER COLUMN value SET DATA TYPE jsonb USING '{}'::jsonb,
    ALTER COLUMN value SET DEFAULT '{}'::jsonb;
-- +goose Down
-- +goose StatementBegin
ALTER TABLE savepoints
    ALTER COLUMN value DROP DEFAULT,
    ALTER COLUMN value SET DATA TYPE TIMESTAMP WITH TIME ZONE,
    ALTER COLUMN value SET DEFAULT '2000-01-01T00:00:00+00:00'::TIMESTAMP WITH TIME ZONE;
-- +goose StatementEnd
