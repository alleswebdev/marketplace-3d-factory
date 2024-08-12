package orders_updater

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

const delayInterval = 10 * time.Second

type OrdersClient interface {
	GetNewOrders(ctx context.Context) (wb.OrdersResponse, error)
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

	err = w.ordersStore.AddOrders(ctx, convertOrders(resp.Orders))
	if err != nil {
		return errors.Wrap(err, "ordersStore.AddOrders")
	}

	log.Println("orders updated")

	return nil
}

func convertOrders(wbOrders []wb.Order) []order.Order {
	result := make([]order.Order, 0, len(wbOrders))
	for _, item := range wbOrders {
		result = append(result, order.Order{
			ID:          item.ID,
			Article:     item.Article,
			Marketplace: card.MpWb.String(),
			OrderCreatedAt: sql.NullTime{
				Time: item.CreatedAt,
			},
		})
	}

	return result
}
