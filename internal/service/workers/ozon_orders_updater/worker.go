package ozon_orders_updater

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
)

const delayInterval = 10 * time.Second

type (
	OrdersClient interface {
		GetUnfulfilledList(ctx context.Context, status string) (ozon.UnfulfilledListResponse, error)
	}
	OrdersStore interface {
		AddOrders(ctx context.Context, orders []orderqueue.Order) error
		GetOrders(ctx context.Context, filter orderqueue.ListFilter) ([]orderqueue.Order, error)
		SetCompleteByOrderIDs(ctx context.Context, orderIDs []string) error
		SetComplete(ctx context.Context, id string, isComplete bool) error
		SetPrinting(ctx context.Context, id string, isPrinting bool) error
		SetChildrenComplete(ctx context.Context, id string, isComplete bool) error
	}
	CardsStore interface {
		AddCards(ctx context.Context, cards []card.Card) error
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
			err := w.update(ctx)
			if err != nil {
				log.Printf("ozon_orders_updater:%s\n", err)
			}
			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) update(ctx context.Context) error {
	resp, err := w.ordersClient.GetUnfulfilledList(ctx, ozon.StatusAwaitingDeliver)
	if err != nil {
		return errors.Wrap(err, "ordersClient.GetUnfulfilledList")
	}

	if len(resp.Result.Postings) <= 0 {
		return nil
	}

	ordersArticles := make([]string, 0, len(resp.Result.Postings))
	for _, item := range resp.Result.Postings {
		for _, product := range item.Products {
			ordersArticles = append(ordersArticles, product.OfferID)
		}
	}

	cards, err := w.cardsStore.GetByArticlesMap(ctx, ordersArticles)
	if err != nil {
		return errors.Wrap(err, "cardsStore.GetByArticlesMap")
	}

	err = w.ordersStore.AddOrders(ctx, convertRespToOrders(resp, cards))
	if err != nil {
		return errors.Wrap(err, "ordersStore.AddOrders")
	}

	//log.Println("ozon orders updated")

	return nil
}

func convertRespToOrders(resp ozon.UnfulfilledListResponse, cards map[string]card.Card) []orderqueue.Order {
	postings := resp.Result.Postings
	result := make([]orderqueue.Order, 0, len(postings))
	for _, item := range postings {
		for _, product := range item.Products {
			c, ok := cards[product.OfferID]
			if !ok {
				continue
			}
			result = append(result, orderqueue.Order{
				ID:             item.PostingNumber,
				Article:        product.OfferID,
				Marketplace:    card.MpOzon.String(),
				Items:          makeItems(c),
				OrderCreatedAt: sql.NullTime{Time: item.InProcessAt, Valid: true},
				Info: orderqueue.Info{
					OrderNumber:     item.PostingNumber,
					OrderShipmentAt: item.ShipmentDate,
					Quantity:        int32(product.Quantity),
				},
			})
		}
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
