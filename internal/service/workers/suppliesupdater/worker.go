// Package suppliesupdater закрывает собранные поставки
package suppliesupdater

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/client/wb"
)

const delayInterval = 5 * time.Second
const StatusDeclinedByClient = "declined_by_client"

type OrdersStore interface {
	GetOrders(ctx context.Context, filter orderqueue.ListFilter) ([]orderqueue.Order, error)
	SetCompleteByOrderIDs(ctx context.Context, orderIDs []string) error
}

type Worker struct {
	wbClient         wb.Client
	ozonClient       ozon.Client
	ordersQueueStore OrdersStore
}

func NewWorker(wbClient wb.Client, ozonClient ozon.Client, ordersQueueStore OrdersStore) Worker {
	return Worker{
		wbClient:         wbClient,
		ozonClient:       ozonClient,
		ordersQueueStore: ordersQueueStore,
	}
}

func (w Worker) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			wbCtxTimeout, wbCancel := context.WithTimeout(ctx, time.Second*30)
			if err := w.updateWb(wbCtxTimeout); err != nil {
				log.Printf("wb_supplies_updater:%s\n", err)
			}

			if err := w.updateWbCancelled(wbCtxTimeout); err != nil {
				log.Printf("wb_supplies_updater_cancelled:%s\n", err)
			}

			wbCancel()

			ozonCtxTimeout, ozonCancel := context.WithTimeout(ctx, time.Second*30)
			if ozonErr := w.updateOzon(ozonCtxTimeout); ozonErr != nil {
				log.Printf("ozon_supplies_updater:%s\n", ozonErr)
			}
			ozonCancel()

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

	orderIDs := make([]string, 0)

	for _, supplyID := range suppliesIDs {
		resp, err := w.wbClient.GetSupplyOrders(ctx, supplyID)
		if err != nil {
			log.Println(errors.Wrap(err, "wbClient.GetSupplyOrders").Error())
			continue
		}

		for _, order := range resp.Orders {
			orderIDs = append(orderIDs, strconv.Itoa(int(order.ID)))
		}
	}

	if err := w.ordersQueueStore.SetCompleteByOrderIDs(ctx, orderIDs); err != nil {
		return errors.Wrap(err, "queueStore.SetCompleteByOrderIDs")
	}

	return nil
}

func (w Worker) updateWbCancelled(ctx context.Context) error {
	orders, err := w.ordersQueueStore.GetOrders(ctx, orderqueue.ListFilter{
		WithParentComplete: false,
		Marketplace:        string(card.MpWb),
	})

	if err != nil {
		return errors.Wrap(err, "ordersQueueStore.GetOrders")
	}

	ids := make([]uint64, 0, len(orders))
	for _, order := range orders {
		id, _ := strconv.ParseUint(order.ID, 10, 64)
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil
	}

	resp, err := w.wbClient.GetOrdersStatus(ctx, ids)
	if err != nil {
		return errors.Wrap(err, "wbClient.GetOrdersStatus")
	}

	var cancelledIDs []string
	for _, order := range resp.Orders {
		if order.SupplierStatus == StatusDeclinedByClient {
			cancelledIDs = append(cancelledIDs, strconv.Itoa(int(order.ID)))
		}
	}

	if len(cancelledIDs) == 0 {
		return nil
	}

	if err := w.ordersQueueStore.SetCompleteByOrderIDs(ctx, cancelledIDs); err != nil {
		return errors.Wrap(err, "SetCompleteByOrderIDs")
	}

	return nil
}

func (w Worker) updateOzon(ctx context.Context) error {
	resp, err := w.ozonClient.GetUnfulfilledList(ctx, ozon.StatusDelivering)
	if err != nil {
		return errors.Wrap(err, "ordersClient.GetUnfulfilledList")
	}

	if len(resp.Result.Postings) <= 0 {
		return nil
	}

	orderIDs := make([]string, 0, len(resp.Result.Postings))
	for _, item := range resp.Result.Postings {
		orderIDs = append(orderIDs, item.PostingNumber)
	}

	if err = w.ordersQueueStore.SetCompleteByOrderIDs(ctx, orderIDs); err != nil {
		return errors.Wrap(err, "ordersQueueStore.SetCompleteByOrderIDs")
	}

	return nil
}
