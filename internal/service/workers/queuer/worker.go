package queuer

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/savepoint"
)

const (
	delayInterval = 10 * time.Second
	limitOrders   = 5
	savepointName = "queuer_worker_last_id"
)

type Worker struct {
	dbPool         *pgxpool.Pool
	savepointStore savepoint.Store
	queueStore     queue.Store
	ordersStore    order.Store
	cardsStore     card.Store
}

func NewWorker(dbPool *pgxpool.Pool, ordersStore order.Store, savepointStore savepoint.Store, queueStore queue.Store, cardsStore card.Store) Worker {
	return Worker{
		dbPool:         dbPool,
		savepointStore: savepointStore,
		queueStore:     queueStore,
		ordersStore:    ordersStore,
		cardsStore:     cardsStore,
	}
}

func (w Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := w.do(ctx)
			if err != nil {
				log.Printf("queuer:%s\n", err)
			}

			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) do(ctx context.Context) error {
	sp, err := w.savepointStore.GetByName(ctx, savepointName)
	if err != nil {
		return errors.Wrap(err, "savepointStore.GetByName")
	}

	orders, err := w.ordersStore.GetLastOrders(ctx, sp.Value.Time, sp.Value.ID, limitOrders)
	if err != nil {
		return errors.Wrap(err, "ordersStore.GetLastOrders")
	}

	ordersArticles := make([]string, 0, len(orders))
	for _, order := range orders {
		ordersArticles = append(ordersArticles, order.Article)
	}

	cards, err := w.cardsStore.GetByArticlesMap(ctx, ordersArticles)
	if err != nil {
		return errors.Wrap(err, "cardsStore.GetByArticlesMap")
	}

	queueItems := make([]queue.Item, 0, len(orders))

	for _, order := range orders {
		card, ok := cards[order.Article]
		if !ok {
			continue
		}

		queueItems = append(queueItems, queue.Item{
			OrderID:        order.ID,
			Article:        order.Article,
			OrderCreatedAt: order.OrderCreatedAt.Time,
		})
		if !card.IsComposite {
			continue
		}
		for _, art := range card.Articles {
			queueItems = append(queueItems, queue.Item{
				OrderID:        order.ID,
				Article:        art,
				OrderCreatedAt: order.OrderCreatedAt.Time,
				Parent:         order.ID,
			})
		}
	}

	if len(queueItems) == 0 {
		return nil
	}

	err = db.TransactionWrapper(ctx, w.dbPool, func(ctx context.Context, txConn db.Conn) error {
		queueTxStore := queue.NewStoreWithTx(w.dbPool)
		err = queueTxStore.AddQueueItems(ctx, queueItems)
		if err != nil {
			return errors.Wrap(err, "queueStore.AddQueueItems")
		}

		lastItem := queueItems[len(queueItems)-1]
		err = w.savepointStore.SetByName(ctx, savepointName, savepoint.Value{
			ID:   lastItem.OrderID,
			Time: lastItem.OrderCreatedAt,
		})
		return errors.Wrap(err, "savepointStore.SetByName")
	})

	if err != nil {
		return errors.Wrap(err, "db.TransactionWrapper")
	}

	log.Println("queue items added")
	return nil
}
