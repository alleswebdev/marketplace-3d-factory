package wbordersupdater

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/wb"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
)

const delayInterval = 10 * time.Second

type (
	OrdersClient interface {
		GetNewOrders(ctx context.Context) (wb.OrdersResponse, error)
	}
	OrdersStore interface {
		AddOrders(ctx context.Context, orders []orderqueue.Order) error
	}
	CardsStore interface {
		GetByArticlesMap(ctx context.Context, articles []string) (map[string]card.Card, error)
	}
)

type Worker struct {
	ordersClient OrdersClient
	ordersStore  OrdersStore
	cardsStore   CardsStore
}

func NewWorker(ordersClient OrdersClient, ordersStore OrdersStore, cardsStore CardsStore) Worker {
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
			if err := w.update(ctx); err != nil {
				log.Printf("wb_orders_updater:%s\n", err)
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

	if err = w.ordersStore.AddOrders(ctx, convertOrders(resp.Orders, cards)); err != nil {
		return errors.Wrap(err, "ordersStore.AddOrders")
	}

	return nil
}

func convertOrders(wbOrders []wb.Order, cards map[string]card.Card) []orderqueue.Order {
	result := make([]orderqueue.Order, 0, len(wbOrders))
	for _, item := range wbOrders {
		c, ok := cards[item.Article]
		if !ok {
			continue
		}

		result = append(result, orderqueue.Order{
			ID:             strconv.Itoa(int(item.ID)),
			Article:        item.Article,
			Items:          makeItems(c),
			Marketplace:    card.MpWb.String(),
			OrderCreatedAt: sql.NullTime{Time: item.CreatedAt},
		})
	}

	return result
}

func makeItems(c card.Card) []orderqueue.Item {
	if !c.IsComposite {
		return []orderqueue.Item{}
	}

	result := make([]orderqueue.Item, 0, len(c.Articles))
	for _, art := range c.Articles {
		result = append(result, orderqueue.Item{
			ID:         uuid.NewString(),
			Name:       art,
			IsComplete: false,
		})
	}

	return result
}
