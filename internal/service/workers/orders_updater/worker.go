package orders_updater

import (
	"context"
	"database/sql"
	"log"
	"time"

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
			w.update(ctx)
			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) update(ctx context.Context) {
	resp, err := w.ordersClient.GetNewOrders(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = w.ordersStore.AddOrders(ctx, convertOrders(resp.Orders))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("orders updated")
}

func convertOrders(wbOrders []wb.Order) []order.Order {
	result := make([]order.Order, 0, len(wbOrders))
	for _, item := range wbOrders {
		result = append(result, order.Order{
			ID:      item.Id,
			Article: item.Article,
			OrderCreatedAt: sql.NullTime{
				Time: item.CreatedAt,
			},
		})
	}

	return result
}
