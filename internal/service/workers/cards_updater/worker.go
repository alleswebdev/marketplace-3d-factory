package cards_updater

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/yandex"

	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

const delayInterval = 5 * time.Second

type (
	CardsStore interface {
		AddCards(ctx context.Context, cards []card.Card) error
		GetByArticlesMap(ctx context.Context, articles []string) (map[string]card.Card, error)
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

			yandexCtxTimeout, yandexCancel := context.WithTimeout(ctx, time.Second*30)
			yandexErr := w.updateYandex(yandexCtxTimeout)
			yandexCancel()
			if err != nil {
				log.Printf("yandex_cards_updater:%s\n", yandexErr)
			}

			time.Sleep(delayInterval)
		}
	}
}

func (w Worker) updateWb(ctx context.Context) error {
	const cardsLimit = 99

	var (
		updatedAt = ""
		nmId      = 0
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

		err = w.cardStore.AddCards(ctx, card.ConvertCards(cardsResp.Cards))
		if err != nil {
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
	cardsResp, err := w.ozonClient.GetProductList(ctx)
	if err != nil {
		return errors.Wrap(err, "ozonClient.GetProductList")
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

	err = w.cardStore.AddCards(ctx, cards)
	return errors.Wrap(err, "cardStore.AddCards")
}
