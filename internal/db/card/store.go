package card

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	tableName         = "cards"
	idColumn          = "id"
	nameColumn        = "name"
	articleColumn     = "article"
	photoColumn       = "photo"
	createdAtColumn   = "created_at"
	updatedAtColumn   = "updated_at"
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

func (s *Store) AddCards(ctx context.Context, cards []Card) error {
	qb := sq.Insert(tableName).
		Columns(idColumn, nameColumn, articleColumn, photoColumn, marketplaceColumn).
		Suffix(
			fmt.Sprintf(`ON CONFLICT(%s, %s) DO NOTHING`, articleColumn, marketplaceColumn),
		).
		PlaceholderFormat(sq.Dollar)

	for _, item := range cards {
		qb = qb.Values(item.ID, item.Name, item.Article, item.Photo, item.Marketplace)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return errors.Wrap(err, "sq.ToSql")
	}

	_, err = s.dbPool.Exec(ctx, query, args...)

	return errors.Wrap(err, "dbPool.Exec")
}

func (s *Store) GetByArticlesMap(ctx context.Context, articles []string) (map[string]Card, error) {
	qb := sq.Select("*").
		From(tableName).
		Where(sq.Eq{articleColumn: articles}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "sq.ToSql")
	}

	var items []Card
	err = pgxscan.Select(ctx, s.dbPool, &items, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "pgxscan.Select")
	}

	byArticlesMap := make(map[string]Card)
	for _, card := range items {
		byArticlesMap[card.Article] = card
	}

	return byArticlesMap, nil
}
