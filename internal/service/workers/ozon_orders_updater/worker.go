package ozon_orders_updater

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
)

const delayInterval = 10 * time.Second

type OrdersClient interface {
	GetUnfulfilledList(ctx context.Context, status string) (ozon.UnfulfilledListResponse, error)
}

type Worker struct {
	ordersClient OrdersClient
	ordersStore  order.Store
}

func NewWorker(ordersClient OrdersClient, ordersStore order.Store) Worker {
	return Worker{
		ordersClient: ordersClient,
		ordersStore:  ordersStore,
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
	resp, err := w.ordersClient.GetUnfulfilledList(ctx, ozon.StatusAwaitingPackaging)
	if err != nil {
		return errors.Wrap(err, "ordersClient.GetUnfulfilledList")
	}

	if len(resp.Result.Postings) <= 0 {
		return nil
	}

	err = w.ordersStore.AddOrders(ctx, convertRespToOrders(resp))
	if err != nil {
		return errors.Wrap(err, "ordersStore.AddOrders")
	}

	log.Println("ozon orders updated")

	return nil
}

func convertRespToOrders(resp ozon.UnfulfilledListResponse) []order.Order {
	postings := resp.Result.Postings
	result := make([]order.Order, 0, len(postings))
	for _, item := range postings {
		for _, product := range item.Products {
			result = append(result, order.Order{
				ID:          item.OrderID,
				Article:     product.OfferID,
				Marketplace: card.MpOzon.String(),
				OrderCreatedAt: sql.NullTime{
					Time:  item.InProcessAt,
					Valid: true,
				},
				OrderShipmentAt: sql.NullTime{
					Time:  item.ShipmentDate,
					Valid: true,
				},
			})
		}
	}

	return result
}
