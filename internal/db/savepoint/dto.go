package savepoint

import (
	"database/sql"
	"time"
)

type Savepoint struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Value     time.Time    `db:"value"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
