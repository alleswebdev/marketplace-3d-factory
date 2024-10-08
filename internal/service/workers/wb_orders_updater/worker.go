package wb_orders_updater

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order_queue"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

const delayInterval = 10 * time.Second

type OrdersClient interface {
	GetNewOrders(ctx context.Context) (wb.OrdersResponse, error)
}

type Worker struct {
	ordersClient OrdersClient
	ordersStore  order_queue.Store
	cardsStore   card.Store
}

func NewWorker(ordersClient OrdersClient, ordersStore order_queue.Store, cardsStore card.Store) Worker {
	return Worker{
		ordersClient: ordersClient,
		ordersStore:  ordersStore,
		cardsStore:   cardsStore,
	}
}

func (w Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := w.update(ctx)
			if err != nil {
				log.Printf("orders_updater:%s\n", err)
			}
			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) update(ctx context.Context) error {
	resp, err := w.ordersClient.GetNewOrders(ctx)
	if err != nil {
		return errors.Wrap(err, "ordersClient.GetNewOrders")
	}

	if len(resp.Orders) <= 0 {
		return nil
	}

	ordersArticles := make([]string, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		ordersArticles = append(ordersArticles, order.Article)
	}

	cards, err := w.cardsStore.GetByArticlesMap(ctx, ordersArticles)
	if err != nil {
		return errors.Wrap(err, "cardsStore.GetByArticlesMap")
	}

	err = w.ordersStore.AddOrders(ctx, convertOrders(resp.Orders, cards))
	if err != nil {
		return errors.Wrap(err, "ordersStore.AddOrders")
	}

	//log.Println("wb orders updated")

	return nil
}

func convertOrders(wbOrders []wb.Order, cards map[string]card.Card) []order_queue.Order {
	result := make([]order_queue.Order, 0, len(wbOrders))
	for _, item := range wbOrders {
		c, ok := cards[item.Article]
		if !ok {
			continue
		}

		result = append(result, order_queue.Order{
			ID:             strconv.Itoa(int(item.ID)),
			Article:        item.Article,
			Items:          makeItems(c),
			Marketplace:    card.MpWb.String(),
			OrderCreatedAt: sql.NullTime{Time: item.CreatedAt},
		})
	}

	return result
}

func makeItems(c card.Card) []order_queue.Item {
	if !c.IsComposite {
		return []order_queue.Item{}
	}

	result := make([]order_queue.Item, 0, len(c.Articles))
	for _, art := range c.Articles {
		result = append(result, order_queue.Item{
			ID:         uuid.NewString(),
			Name:       art,
			IsComplete: false,
		})
	}

	return result
}
