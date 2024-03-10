package order

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	tableName = "orders"

	idColumn                = "id"
	articleColumn           = "article"
	orderCreatedAtColumn    = "order_created_at"
	orderShipmentDateColumn = "order_shipment_date"
	createdAtColumn         = "created_at"
	updatedAtColumn         = "updated_at"
	marketplaceColumn       = "marketplace"
)

type Store struct {
	dbPool *pgxpool.Pool
}

func New(dbPool *pgxpool.Pool) Store {
	return Store{dbPool: dbPool}
}

func (s *Store) AddOrder(ctx context.Context, order Order) error {
	qb := sq.Insert(tableName).
		Columns(idColumn, articleColumn, orderCreatedAtColumn, orderShipmentDateColumn, marketplaceColumn).
		Values(order.ID, order.Article, order.OrderCreatedAt, order.OrderShipmentAt, order.Marketplace).
		Suffix(
			fmt.Sprintf(`ON CONFLICT(%s) DO NOTHING`, idColumn),
		).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) AddOrders(ctx context.Context, orders []Order) error {
	qb := sq.Insert(tableName).
		Columns(idColumn, articleColumn, orderCreatedAtColumn, orderShipmentDateColumn, marketplaceColumn).
		Suffix(
			fmt.Sprintf(`ON CONFLICT(%s) DO NOTHING`, idColumn),
		).
		PlaceholderFormat(sq.Dollar)

	for _, item := range orders {
		qb = qb.Values(item.ID, item.Article, item.OrderCreatedAt.Time, item.OrderShipmentAt, item.Marketplace)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) GetLastOrders(ctx context.Context, marketplace string, lastCreatedAt time.Time, lastID int64, limit int64) ([]Order, error) {
	qb := sq.Select("*").
		From(tableName).
		Limit(uint64(limit)).
		OrderBy(orderCreatedAtColumn, idColumn).
		Where(sq.Eq{marketplaceColumn: marketplace}).
		PlaceholderFormat(sq.Dollar)

	if lastID != 0 {
		qb = qb.Where(
			sq.Or{
				sq.Gt{orderCreatedAtColumn: lastCreatedAt},
				sq.And{
					sq.Eq{orderCreatedAtColumn: lastCreatedAt},
					sq.Gt{idColumn: lastID},
				},
			},
		)
	} else {
		qb = qb.Where(sq.GtOrEq{orderCreatedAtColumn: lastCreatedAt})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sq.ToSql")
	}

	var items []Order
	err = pgxscan.Select(ctx, s.dbPool, &items, query, args...)

	return items, errors.Wrap(err, "pgxscan.Select")
}
