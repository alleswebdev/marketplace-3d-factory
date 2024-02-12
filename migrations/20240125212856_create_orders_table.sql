-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION upd_updated_at() RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION add_time_fields(table_name text) RETURNS VOID
    LANGUAGE plpgsql
AS $$
DECLARE
    trigger_name TEXT;
BEGIN
    EXECUTE 'ALTER TABLE ' || table_name || ' ADD COLUMN created_at timestamp WITH TIME ZONE DEFAULT NOW() NOT NULL;';

    EXECUTE 'ALTER TABLE ' || table_name || ' ADD COLUMN updated_at timestamp WITH TIME ZONE DEFAULT NOW() NOT NULL;';

    trigger_name := 't_' || replace(table_name, '.', '_') || '_upt';
    EXECUTE 'CREATE TRIGGER ' || trigger_name || ' BEFORE UPDATE ON ' || table_name ||
            ' FOR EACH ROW EXECUTE PROCEDURE upd_updated_at()';
END;
$$;
-- +goose StatementEnd

CREATE TABLE orders (
                        id                 int8             NOT NULL,
                        article            text             NOT NULL,
                        order_created_at   timestamptz      NOT NULL,
                        CONSTRAINT orders_pkey PRIMARY KEY (id)
);

SELECT add_time_fields('orders');

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
