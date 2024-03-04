-- +goose Up
ALTER TABLE orders
   ADD COLUMN   order_shipment_date timestamptz      DEFAULT NULL,
   ADD COLUMN   marketplace TEXT DEFAULT 'wb';

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN   order_shipment_date,
    DROP COLUMN   marketplace;
-- +goose StatementEnd
