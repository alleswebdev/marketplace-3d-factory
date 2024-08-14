package cards_updater

import (
	"context"
	"log"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

const delayInterval = 5 * time.Minute

type Worker struct {
	wbClient   wb.Client
	ozonClient ozon.Client
	cardStore  card.Store
}

func NewWorker(wbClient wb.Client, ozonClient ozon.Client, cardStore card.Store) Worker {
	return Worker{
		wbClient:   wbClient,
		ozonClient: ozonClient,
		cardStore:  cardStore,
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
				log.Printf("wb_cards_updater:%s\n", err)
			}

			ozonCtxTimeout, ozonCancel := context.WithTimeout(ctx, time.Second*30)
			ozonErr := w.updateOzon(ozonCtxTimeout)
			ozonCancel()
			if err != nil {
				log.Printf("ozon_cards_updater:%s\n", ozonErr)
			}

			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) updateWb(ctx context.Context) error {
	cardsResp, err := w.wbClient.GetCardsList(ctx)
	if err != nil {
		return errors.Wrap(err, "wbClient.GetCardsList")
	}

	err = w.cardStore.AddCards(ctx, card.ConvertCards(cardsResp.Cards))
	if err != nil {
		return errors.Wrap(err, "cardStore.AddCards")
	}

	return nil
}

func (w Worker) updateOzon(ctx context.Context) error {
	cardsResp, err := w.ozonClient.GetProductList(ctx)
	if err != nil {
		errors.Wrap(err, "ozonClient.GetProductList")
	}

	productIDs := make([]int64, 0, len(cardsResp.Result.Items))
	for _, item := range cardsResp.Result.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := w.ozonClient.GetProductInfoList(ctx, productIDs)
	if err != nil {
		return errors.Wrap(err, "ozonClient.GetProductInfoList")
	}

	err = w.cardStore.AddCards(ctx, convertProductResponseToCards(products))
	return errors.Wrap(err, "cardStore.AddCards")
}

func convertProductResponseToCards(productsResponse ozon.ProductListInfoResponse) []card.Card {
	result := make([]card.Card, 0, len(productsResponse.Result.Items))
	for _, item := range productsResponse.Result.Items {
		convertItem := card.Card{
			ID:          uuid.New(),
			Name:        item.Name,
			Article:     item.OfferID,
			Marketplace: card.MpOzon,
			IsComposite: false,
			Photo:       item.PrimaryImage,
		}

		result = append(result, convertItem)
	}

	return result
}
