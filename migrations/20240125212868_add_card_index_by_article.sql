-- +goose Up
CREATE UNIQUE INDEX cards_article_marketplace on cards(article, marketplace);

-- +goose Down
-- +goose StatementBegin
DROP INDEX cards_article_marketplace;
-- +goose StatementEnd
