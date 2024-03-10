package order

import (
	"database/sql"
)

type Order struct {
	ID              int64        `db:"id"`
	Article         string       `db:"article"`
	Marketplace     string       `db:"marketplace"`
	OrderCreatedAt  sql.NullTime `db:"order_created_at"`
	OrderShipmentAt sql.NullTime `db:"order_shipment_date"`
	CreatedAt       sql.NullTime `db:"created_at"`
	UpdatedAt       sql.NullTime `db:"updated_at"`
}
