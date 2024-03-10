-- +goose Up
ALTER TABLE queue
   ADD COLUMN   order_shipment_date timestamptz      DEFAULT NULL,
   ADD COLUMN   marketplace TEXT DEFAULT 'wb';

-- +goose Down
-- +goose StatementBegin
ALTER TABLE queue
    DROP COLUMN   order_shipment_date,
    DROP COLUMN   marketplace;
-- +goose StatementEnd
