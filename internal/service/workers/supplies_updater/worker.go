package supplies_updater

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

const delayInterval = 5 * time.Second

type Worker struct {
	wbClient    wb.Client
	ozonClient  ozon.Client
	ordersStore order.Store
	queueStore  queue.Store
}

func NewWorker(wbClient wb.Client, ozonClient ozon.Client, ordersStore order.Store, queueStore queue.Store) Worker {
	return Worker{
		wbClient:    wbClient,
		ordersStore: ordersStore,
		queueStore:  queueStore,
		ozonClient:  ozonClient,
	}
}

func (w Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			wbCtxTimeout, wbCancel := context.WithTimeout(ctx, time.Second*30)
			err := w.updateWb(wbCtxTimeout)
			wbCancel()
			if err != nil {
				log.Printf("wb_supplies_updater:%s\n", err)
			}

			ozonCtxTimeout, ozonCancel := context.WithTimeout(ctx, time.Second*30)
			ozonErr := w.updateOzon(ozonCtxTimeout)
			ozonCancel()
			if err != nil {
				log.Printf("ozon_supplies_updater:%s\n", ozonErr)
			}

			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) updateWb(ctx context.Context) error {
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

func (w Worker) updateOzon(ctx context.Context) error {
	resp, err := w.ozonClient.GetUnfulfilledList(ctx, ozon.StatusAwaitingDeliver)
	if err != nil {
		return errors.Wrap(err, "ordersClient.GetUnfulfilledList")
	}

	if len(resp.Result.Postings) <= 0 {
		return nil
	}

	orderIDs := make([]int64, 0, len(resp.Result.Postings))
	for _, item := range resp.Result.Postings {
		orderIDs = append(orderIDs, item.OrderID)
	}

	err = w.queueStore.SetCompleteByOrderIDs(ctx, orderIDs)
	return errors.Wrap(err, "queueStore.SetCompleteByOrderIDs")
}
