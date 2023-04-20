-- +goose Up
CREATE TABLE cards (
                       id      UUID PRIMARY KEY,
                       name    TEXT NOT NULL,
                       article TEXT NOT NULL,
                       photo   TEXT DEFAULT NULL,
                       articles TEXT[] NOT NULL,
                       color TEXT DEFAULT 'standart',
                       size TEXT DEFAULT 'standart',
                       marketplace TEXT DEFAULT 'wb',
                       is_composite BOOL DEFAULT false
);

SELECT add_time_fields('cards');
-- +goose Down
-- +goose StatementBegin
drop table if exists cards;
-- +goose StatementEnd
