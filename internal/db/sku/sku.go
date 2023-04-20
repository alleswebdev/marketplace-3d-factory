package sku

import (
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	tableName = "skus"

	idColumn          = "id"
	nmIDColumn        = "nmID"
	nameColumn        = "name"
	articlesColumn    = "articles"
	filesColumn       = "files"
	colorColumn       = "color"
	sizeColumn        = "size"
	marketplaceColumn = "marketplace"
	isCompositeColumn = "is_composite"
)

type Store struct {
	dbPool *pgxpool.Pool
}

func New(dbPool *pgxpool.Pool) Store {
	return Store{dbPool: dbPool}
}

func (s *Store) AddSKUs(ctx context.Context, SKUs []SKU) error {
	qb := sq.Insert(tableName).
		Columns(nameColumn, articlesColumn, colorColumn, sizeColumn, marketplaceColumn, isCompositeColumn).
		Suffix(
			fmt.Sprintf(`ON CONFLICT(%s) DO NOTHING`, idColumn),
		).
		PlaceholderFormat(sq.Dollar)

	for _, item := range SKUs {
		qb = qb.Values(item.Name, item.Articles, item.Color, item.Size, item.Marketplace, item.IsComposite)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) SetNmIDs(ctx context.Context, SKUs []SKU) error {
	var values []string
	for _, data := range SKUs {
		values = append(values, fmt.Sprintf("('%s'::UUID, '%s')", data.NmID, data.Name))
	}

	query := fmt.Sprintf(`
		UPDATE skus
		SET nmID = data.nmID
		FROM (VALUES %s) AS data (nmID, name)
		WHERE skus.name = data.name;
	`, strings.Join(values, ","))

	_, err := s.dbPool.Exec(ctx, query)

	return errors.Wrap(err, "dbPool.Exec")
}
