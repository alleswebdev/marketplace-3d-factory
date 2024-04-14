-- +goose Up
CREATE UNIQUE INDEX orders_article_id on orders(article, id);

-- +goose Down
-- +goose StatementBegin
DROP INDEX orders_article_id;
-- +goose StatementEnd
