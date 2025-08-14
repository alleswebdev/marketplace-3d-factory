package yandexordersupdater

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/yandex"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const delayInterval = 10 * time.Second

type (
	OrdersClient interface {
		GetOrders(ctx context.Context, status string) (yandex.OrdersDTO, error)
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
				log.Printf("yandexordersupdater:%s\n", err)
			}

			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) update(ctx context.Context) error {
	resp, err := w.ordersClient.GetOrders(ctx, "PROCESSING")
	if err != nil {
		return errors.Wrap(err, "ordersClient.GetOrders")
	}

	if len(resp.Orders) <= 0 {
		return nil
	}

	ordersArticles := make([]string, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		for _, product := range order.Items {
			ordersArticles = append(ordersArticles, product.OfferId)
		}
	}

	cards, err := w.cardsStore.GetByArticlesMap(ctx, ordersArticles)
	if err != nil {
		return errors.Wrap(err, "cardsStore.GetByArticlesMap")
	}

	if err = w.ordersStore.AddOrders(ctx, convertRespToOrders(resp, cards)); err != nil {
		return errors.Wrap(err, "ordersStore.AddOrders")
	}

	return nil
}

func convertRespToOrders(resp yandex.OrdersDTO, cards map[string]card.Card) []orderqueue.Order {
	result := make([]orderqueue.Order, 0, len(resp.Orders))
	for _, order := range resp.Orders {
		if order.Substatus == "SHIPPED" {
			continue
		}
		for _, product := range order.Items {
			c, ok := cards[product.OfferId]
			if !ok {
				continue
			}

			createdAt, err := time.ParseInLocation("02-01-2006 15:04:05", order.CreationDate, time.FixedZone("MSK", 3*60*60))
			if err != nil {
				createdAt = time.Now()
			}

			shipmentAt := time.Now()
			if len(order.Delivery.Shipments) > 0 {
				shipmentAt, err = time.Parse("02-01-2006", order.Delivery.Shipments[0].ShipmentDate)
				if err != nil {
					shipmentAt = time.Now()
				}
			}

			result = append(result, orderqueue.Order{
				ID:             strconv.Itoa(product.Id),
				Article:        product.OfferId,
				Marketplace:    card.MpYandex.String(),
				Items:          makeItems(c),
				OrderCreatedAt: sql.NullTime{Time: createdAt, Valid: true},
				Info: orderqueue.Info{
					OrderNumber:     fmt.Sprintf("№ %[1]d / %[1]d", order.Id),
					OrderShipmentAt: shipmentAt,
					Quantity:        int32(product.Count),
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
