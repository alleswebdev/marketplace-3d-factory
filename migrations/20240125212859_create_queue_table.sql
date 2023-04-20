-- +goose Up
CREATE TABLE queue (
                       id serial PRIMARY KEY,
                       order_id BIGINT NOT NULL,
                       article TEXT NOT NULL,
                       parent BIGINT NOT NULL DEFAULT 0,
                       is_printing BOOL DEFAULT false,
                       is_complete BOOL DEFAULT false,
                       order_created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

SELECT add_time_fields('queue');
-- +goose Down
-- +goose StatementBegin
drop table if exists queue;
-- +goose StatementEnd
