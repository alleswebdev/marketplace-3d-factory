package supplies_updater

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

const delayInterval = 5 * time.Second

type Worker struct {
	wbClient    wb.Client
	ordersStore order.Store
	queueStore  queue.Store
}

func NewWorker(wbClient wb.Client, ordersStore order.Store, queueStore queue.Store) Worker {
	return Worker{
		wbClient:    wbClient,
		ordersStore: ordersStore,
		queueStore:  queueStore,
	}
}

func (w Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*30)
			err := w.update(ctxTimeout)
			cancel()
			if err != nil {
				log.Printf("supplies_updater:%s\n", err)
			}
			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) update(ctx context.Context) error {
	next := 1
	suppliesIDs := make([]string, 0)

	for {
		if next <= 0 {
			break
		}

		resp, err := w.wbClient.GetSupplies(ctx, next)
		if err != nil {
			return errors.Wrap(err, "wbClient.GetSupplies")
		}

		if len(resp.Supplies) <= 0 {
			break
		}
		next = resp.Next

		for _, supply := range resp.Supplies {
			if supply.Done {
				continue
			}
			suppliesIDs = append(suppliesIDs, supply.ID)
		}
	}

	orderIDs := make([]int64, 0)

	for _, supplyID := range suppliesIDs {
		//todo расспаралелить тут запросы, если вб не ограничивает лимиты на эту ручку
		resp, err := w.wbClient.GetSupplyOrders(ctx, supplyID)
		if err != nil {
			log.Println(errors.Wrap(err, "wbClient.GetSupplyOrders").Error())
			continue
		}

		for _, order := range resp.Orders {
			orderIDs = append(orderIDs, order.ID)
		}
	}

	err := w.queueStore.SetCompleteByOrderIDs(ctx, orderIDs)
	return errors.Wrap(err, "queueStore.SetCompleteByOrderIDs")
}
