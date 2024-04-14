package queue

import (
	"database/sql"
	"time"
)

type Item struct {
	ID              int64        `db:"id"`
	OrderID         int64        `db:"order_id"`
	Article         string       `db:"article"`
	Parent          int64        `db:"parent"`
	IsPrinting      bool         `db:"is_printing"`
	IsComplete      bool         `db:"is_complete"`
	Marketplace     string       `db:"marketplace"`
	OrderCreatedAt  time.Time    `db:"order_created_at"`
	OrderShipmentAt time.Time    `db:"order_shipment_date"`
	CreatedAt       sql.NullTime `db:"created_at"`
	UpdatedAt       sql.NullTime `db:"updated_at"`
	Info            Info         `db:"info"`
}

type Info struct {
	OrderNumber     string    `json:"order_number"`
	OrderShipmentAt time.Time `json:"order_shipment_date"`
	Quantity        int32     `json:"quantity"`
}
