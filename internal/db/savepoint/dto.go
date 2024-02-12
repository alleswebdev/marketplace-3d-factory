package savepoint

import (
	"database/sql"
	"time"
)

type Value struct {
	ID   int64     `db:"id"`
	Time time.Time `db:"value"`
}

type Savepoint struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Value     Value        `db:"value"`
	CreatedAt sql.NullTime `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
