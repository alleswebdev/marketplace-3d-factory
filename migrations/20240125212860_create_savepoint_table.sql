-- +goose Up
CREATE TABLE savepoints (
                       id      serial PRIMARY KEY,
                       name TEXT NOT NULL unique ,
                       value     TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT '2000-01-01T00:00:00+00:00'::TIMESTAMP WITH TIME ZONE
);

SELECT add_time_fields('savepoints');
-- +goose Down
-- +goose StatementBegin
drop table if exists savepoints;
-- +goose StatementEnd
