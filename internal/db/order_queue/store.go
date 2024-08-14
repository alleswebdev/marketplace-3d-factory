package order_queue

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	tableName = "orders_queue"

	idColumn             = "id"
	articleColumn        = "article"
	itemsColumn          = "order_composite_items"
	orderCreatedAtColumn = "order_created_at"
	createdAtColumn      = "created_at"
	updatedAtColumn      = "updated_at"
	marketplaceColumn    = "marketplace"
	infoColumn           = "info"
	isCompleteColumn     = "is_complete"
	isPrintingColumn     = "is_printing"
)

type Store struct {
	dbPool *pgxpool.Pool
}

func New(dbPool *pgxpool.Pool) Store {
	return Store{dbPool: dbPool}
}

func (s *Store) AddOrders(ctx context.Context, orders []Order) error {
	qb := sq.Insert(tableName).
		Columns(idColumn, articleColumn, orderCreatedAtColumn, itemsColumn, marketplaceColumn, infoColumn).
		Suffix(
			fmt.Sprintf(`ON CONFLICT(%s, %s) DO NOTHING`, articleColumn, idColumn),
		).
		PlaceholderFormat(sq.Dollar)

	for _, item := range orders {
		qb = qb.Values(item.ID, item.Article, item.OrderCreatedAt.Time, item.Items, item.Marketplace, item.Info)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) GetOrders(ctx context.Context, filter ListFilter) ([]Order, error) {
	qb := sq.Select("*").
		From(tableName).
		Limit(100).
		Where(sq.Eq{marketplaceColumn: filter.GetMarketplace()}).
		Where(sq.Gt{createdAtColumn: time.Now().Add(-time.Hour * 24 * 7)}).
		Where(sq.Eq{isCompleteColumn: filter.WithParentComplete}).
		PlaceholderFormat(sq.Dollar)

	if filter.GetMarketplace() == card.MpOzon.String() {
		qb = qb.OrderBy(`info->>'order_shipment_date'`, orderCreatedAtColumn)
	} else {
		qb = qb.OrderBy(orderCreatedAtColumn)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sq.ToSql")
	}

	var items []Order
	err = pgxscan.Select(ctx, s.dbPool, &items, query, args...)

	return items, errors.Wrap(err, "pgxscan.Select")
}

func (s *Store) SetCompleteByOrderIDs(ctx context.Context, orderIDs []string) error {
	qb := sq.Update(tableName).
		Set(isCompleteColumn, true).
		Where(sq.Eq{idColumn: orderIDs}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) SetComplete(ctx context.Context, id string, isComplete bool) error {
	qb := sq.Update(tableName).
		Set(isCompleteColumn, isComplete).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) SetPrinting(ctx context.Context, id string, isPrinting bool) error {
	qb := sq.Update(tableName).
		Set(isPrintingColumn, isPrinting).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) SetChildrenComplete(ctx context.Context, id string, isComplete bool) error {
	qb := sq.Update(tableName).
		Set(itemsColumn, sq.Expr(`
            (
                SELECT jsonb_agg(
                    CASE
                        WHEN item->>'id' = ? THEN item || '{"is_complete": `+strconv.FormatBool(isComplete)+`}'::jsonb
                        ELSE item
                    END
                )
                FROM jsonb_array_elements(order_composite_items) AS item
            )`, id)).
		Where(`order_composite_items @> ?::jsonb`, `[{"id": `+strconv.Quote(id)+`}]`).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}
