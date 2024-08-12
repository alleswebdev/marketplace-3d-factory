-- +goose Up
ALTER TABLE cards DROP CONSTRAINT cards_pkey;

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards ADD PRIMARY KEY (id);
-- +goose StatementEnd
