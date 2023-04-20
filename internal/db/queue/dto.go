package queue

import (
	"database/sql"
	"time"
)

type QueueItem struct {
	ID             int64        `db:"id"`
	OrderID        int64        `db:"order_id"`
	Article        string       `db:"article"`
	Parent         int64        `db:"parent"`
	IsPrinting     bool         `db:"is_printing"`
	IsComplete     bool         `db:"is_complete"`
	OrderCreatedAt time.Time    `db:"order_created_at"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
}
