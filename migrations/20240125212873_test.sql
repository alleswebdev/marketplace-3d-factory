-- +goose Up
CREATE TABLE test (
           id text NOT NULL,
           article text NOT NULL
);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS test;
-- +goose StatementEnd
