package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type Conn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func TransactionWrapper(
	ctx context.Context,
	conn Conn,
	wrappedFunc func(ctx context.Context, txConn Conn) error,
) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "conn.Begin")
	}

	err = wrappedFunc(ctx, tx)

	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Wrapf(err, "rollback error: %s", rollbackErr.Error())
		}

		return errors.Wrap(err, "wrapForRecoverPanic")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "tx.Commit")
	}

	return nil
}
