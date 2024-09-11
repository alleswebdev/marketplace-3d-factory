package cards_updater

import (
	"github.com/google/uuid"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
)

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
