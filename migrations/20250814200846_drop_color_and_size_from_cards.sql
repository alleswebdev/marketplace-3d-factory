-- +goose Up
ALTER TABLE cards
    DROP COLUMN color,
    DROP COLUMN size;

-- +goose Down
ALTER TABLE cards
    ADD COLUMN color text DEFAULT 'standart',
    ADD COLUMN size text DEFAULT 'standart';