package savepoint

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	tableName   = "savepoints"
	idColumn    = "id"
	nameColumn  = "name"
	valueColumn = "value"
)

type Store struct {
	dbPool *pgxpool.Pool
}

func New(dbPool *pgxpool.Pool) Store {
	return Store{dbPool: dbPool}
}

func (s *Store) SetByName(ctx context.Context, name string, value Value) error {
	qb := sq.Insert(tableName).
		Columns(nameColumn, valueColumn).
		Values(name, value).
		Suffix(
			fmt.Sprintf(
				`ON CONFLICT(%s) DO UPDATE SET %s = excluded.%s`,
				nameColumn, valueColumn, valueColumn,
			),
		).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) GetByName(ctx context.Context, name string) (Savepoint, error) {
	qb := sq.Select("*").
		From(tableName).
		Where(sq.Eq{nameColumn: name}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return Savepoint{}, errors.Wrap(err, "sq.ToSql")
	}

	var item Savepoint
	err = pgxscan.Get(ctx, s.dbPool, &item, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return Savepoint{}, nil
	}

	return item, errors.Wrap(err, "pgxscan.Select")
}
