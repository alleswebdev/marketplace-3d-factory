package cardsupdater

import (
	"context"
	"log"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/client/wb"
	"github.com/alleswebdev/marketplace-3d-factory/internal/client/yandex"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const delayInterval = 5 * time.Second

type (
	CardsStore interface {
		AddCards(ctx context.Context, cards []card.Card) error
	}
)

type Worker struct {
	wbClient     wb.Client
	ozonClient   ozon.Client
	yandexClient yandex.Client
	cardStore    CardsStore
}

func NewWorker(wbClient wb.Client, ozonClient ozon.Client, yandexClient yandex.Client, cardStore CardsStore) Worker {
	return Worker{
		wbClient:     wbClient,
		ozonClient:   ozonClient,
		cardStore:    cardStore,
		yandexClient: yandexClient,
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
				log.Printf("wb_cards_updater:%s\n", err)
			}

			wbCancel()

			ozonCtxTimeout, ozonCancel := context.WithTimeout(ctx, time.Second*30)
			if err := w.updateOzon(ozonCtxTimeout); err != nil {
				log.Printf("ozon_cards_updater:%s\n", err)
			}
			ozonCancel()

			yandexCtxTimeout, yandexCancel := context.WithTimeout(ctx, time.Second*30)
			if err := w.updateYandex(yandexCtxTimeout); err != nil {
				log.Printf("yandex_cards_updater:%s\n", err)
			}
			yandexCancel()

			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) updateWb(ctx context.Context) error {
	var (
		updatedAt  = ""
		nmId       = 0
		cardsLimit = 99
	)

	for {
		cardsResp, err := w.wbClient.GetCardsList(ctx, wb.CardListCursor{
			UpdatedAt: updatedAt,
			NmID:      nmId,
			Limit:     cardsLimit,
		})

		if err != nil {
			return errors.Wrap(err, "wbClient.GetCardsList")
		}

		if err = w.cardStore.AddCards(ctx, card.ConvertCards(cardsResp.Cards)); err != nil {
			return errors.Wrap(err, "cardStore.AddCards")
		}

		if cardsResp.CardsListResponseCursor.Total < cardsLimit {
			break
		}

		updatedAt = cardsResp.CardsListResponseCursor.UpdatedAt
		nmId = cardsResp.CardsListResponseCursor.NmID
	}

	return nil
}

func (w Worker) updateOzon(ctx context.Context) error {
	var (
		lastID string
		limit  = 300
	)

	for {
		cardsResp, err := w.ozonClient.GetProductList(ctx, lastID, limit)
		if err != nil {
			return errors.Wrap(err, "ozonClient.GetProductList")
		}

		productIDs := make([]int64, 0, len(cardsResp.Result.Items))
		for _, item := range cardsResp.Result.Items {
			productIDs = append(productIDs, item.ProductID)
		}

		if len(productIDs) == 0 {
			break
		}

		products, err := w.ozonClient.GetProductInfoList(ctx, productIDs)
		if err != nil {
			return errors.Wrap(err, "ozonClient.GetProductInfoList")
		}

		if err = w.cardStore.AddCards(ctx, convertProductResponseToCards(products)); err != nil {
			return errors.Wrap(err, "cardStore.AddCards")
		}

		lastID = cardsResp.Result.LastID
		if len(cardsResp.Result.Items) < limit || cardsResp.Result.LastID == "" {
			break
		}
	}

	return nil
}

func (w Worker) updateYandex(ctx context.Context) error {
	productsResp, err := w.yandexClient.GetProductList(ctx)
	if err != nil {
		return errors.Wrap(err, "yandexClient.GetProductList")
	}

	cards := make([]card.Card, 0, len(productsResp.Result.OfferMappings))
	for _, offerMappings := range productsResp.Result.OfferMappings {
		var photo string
		if len(offerMappings.Offer.Pictures) > 0 {
			photo = offerMappings.Offer.Pictures[0]
		}
		cards = append(cards, card.Card{
			ID:          uuid.New(),
			Name:        offerMappings.Offer.Name,
			Article:     offerMappings.Offer.OfferId,
			Marketplace: card.MpYandex,
			Photo:       photo,
		})
	}

	if err = w.cardStore.AddCards(ctx, cards); err != nil {
		return errors.Wrap(err, "cardStore.AddCards")
	}

	return nil
}
