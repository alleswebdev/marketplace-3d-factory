package order

import (
	"database/sql"
)

type Order struct {
	ID             int64        `db:"id"`
	Article        string       `db:"article"`
	OrderCreatedAt sql.NullTime `db:"order_created_at"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
}
