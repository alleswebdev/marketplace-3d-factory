-- +goose Up
CREATE TABLE public.orders_queue (
                               id text NOT NULL,
                               article text NOT NULL,
                               order_composite_items jsonb DEFAULT '{}'::jsonb,
                               order_created_at timestamptz NOT NULL,
                               created_at timestamptz DEFAULT now() NOT NULL,
                               updated_at timestamptz DEFAULT now() NOT NULL,
                               marketplace text DEFAULT 'wb'::text NULL,
                               info jsonb DEFAULT '{}'::jsonb,
                               is_complete bool default false,
                               is_printing bool default false

);
CREATE UNIQUE INDEX orders_queue_article_id ON public.orders_queue (article,id);

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders_queue;
-- +goose StatementEnd
