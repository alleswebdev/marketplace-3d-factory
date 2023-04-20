package queue

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db"
)

const (
	TableName        = "queue"
	IDColumn         = "id"
	OrderIDColumn    = "order_id"
	ArticleColumn    = "article"
	ParentColumn     = "parent"
	IsCompleteColumn = "is_complete"
	IsPrintingColumn = "is_printing"
	CreatedAtColumn  = "order_created_at"
)

type Store struct {
	dbPool db.Conn
}

func New(dbPool db.Conn) Store {
	return Store{dbPool: dbPool}
}

func NewStoreWithTx(txConn db.Conn) *Store {
	return &Store{
		dbPool: txConn,
	}
}

func (s *Store) AddQueueItems(ctx context.Context, items []QueueItem) error {
	qb := sq.Insert(TableName).
		Columns(OrderIDColumn, ArticleColumn, ParentColumn, CreatedAtColumn).
		PlaceholderFormat(sq.Dollar)

	for _, item := range items {
		qb = qb.Values(item.OrderID, item.Article, item.Parent, item.OrderCreatedAt)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) GetAllItems(ctx context.Context) ([]QueueItem, error) {
	qb := sq.Select("*").
		From(TableName).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sq.ToSql")
	}

	var items []QueueItem
	err = pgxscan.Select(ctx, s.dbPool, &items, query, args...)

	return items, errors.Wrap(err, "pgxscan.Select")
}

func (s *Store) SetComplete(ctx context.Context, id int64, isComplete bool) error {
	qb := sq.Update(TableName).
		Set(IsCompleteColumn, isComplete).
		Where(sq.Eq{IDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) SetPrinting(ctx context.Context, id int64, isPrinting bool) error {
	qb := sq.Update(TableName).
		Set(IsPrintingColumn, isPrinting).
		Where(sq.Eq{IDColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}
