package queue

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
)

const (
	TableName            = "queue"
	IDColumn             = "id"
	OrderIDColumn        = "order_id"
	ArticleColumn        = "article"
	ParentColumn         = "parent"
	IsCompleteColumn     = "is_complete"
	IsPrintingColumn     = "is_printing"
	OrderCreatedAtColumn = "order_created_at"
	OrderShipmentColumn  = "order_shipment_date"
	MarketplaceColumn    = "marketplace"
	InfoColumn           = "info"
)

type Store struct {
	dbPool      db.Conn
	marketplace string
}

func (s *Store) SetMarketplace(marketplace string) {
	s.marketplace = marketplace
}
func (s *Store) GetMarketplace() string {
	if len(s.marketplace) == 0 {
		return card.MpWb.String()
	}

	return s.marketplace
}

func New(dbPool db.Conn) Store {
	return Store{dbPool: dbPool}
}

func NewStoreWithTx(txConn db.Conn) *Store {
	return &Store{
		dbPool: txConn,
	}
}

func (s *Store) AddQueueItems(ctx context.Context, items []Item) error {
	qb := sq.Insert(TableName).
		Columns(OrderIDColumn, ArticleColumn, ParentColumn, OrderCreatedAtColumn, OrderShipmentColumn, MarketplaceColumn, InfoColumn).
		PlaceholderFormat(sq.Dollar)

	for _, item := range items {
		qb = qb.Values(item.OrderID, item.Article, item.Parent, item.OrderCreatedAt, item.OrderShipmentAt, item.Marketplace, item.Info)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) GetAllItems(ctx context.Context) ([]Item, error) {
	qb := sq.Select("*").
		From(TableName).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sq.ToSql")
	}

	var items []Item
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

type ListFilter struct {
	WithParentComplete   bool   `json:"withParentComplete"`
	WithChildrenComplete bool   `json:"withChildrenComplete"`
	Marketplace          string `json:"marketplace"`
}

func (f ListFilter) GetMarketplace() string {
	if len(f.Marketplace) == 0 {
		return card.MpWb.String()
	}

	return f.Marketplace
}

func (s *Store) GetList(ctx context.Context, filter ListFilter) ([]Item, error) {
	qb := sq.Select("*").
		From(TableName).
		PlaceholderFormat(sq.Dollar)

	if filter.GetMarketplace() == card.MpOzon.String() {
		qb = qb.OrderBy(OrderShipmentColumn, OrderCreatedAtColumn)
	} else {
		qb = qb.OrderBy(OrderCreatedAtColumn)
	}

	qb = qb.Where(sq.Eq{MarketplaceColumn: filter.GetMarketplace()})

	wheres := sq.Or{}

	//if !filter.WithChildrenComplete {
	//	wheres = append(wheres, sq.And{
	//		sq.NotEq{ParentColumn: 0},
	//		sq.Eq{IsCompleteColumn: false},
	//	})
	//} else {
	//	wheres = append(wheres, sq.And{
	//		sq.NotEq{ParentColumn: 0},
	//		sq.Eq{IsCompleteColumn: true},
	//	})
	//}

	wheres = append(wheres, sq.And{
		sq.Eq{ParentColumn: 0},
		sq.Eq{IsCompleteColumn: filter.WithParentComplete},
	})
	wheres = append(wheres, sq.And{
		sq.NotEq{ParentColumn: 0},
	})

	if filter.WithParentComplete {
		qb = qb.Where(sq.Gt{OrderCreatedAtColumn: time.Now().Add(-time.Hour * 24 * 7)})
	}

	if len(wheres) > 0 {
		qb = qb.Where(wheres)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sq.ToSql")
	}

	var items []Item
	err = pgxscan.Select(ctx, s.dbPool, &items, query, args...)

	return items, errors.Wrap(err, "pgxscan.Select")
}

func (s *Store) SetCompleteByOrderIDs(ctx context.Context, orderIDs []int64) error {
	qb := sq.Update(TableName).
		Set(IsCompleteColumn, true).
		Where(sq.Eq{OrderIDColumn: orderIDs}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}
