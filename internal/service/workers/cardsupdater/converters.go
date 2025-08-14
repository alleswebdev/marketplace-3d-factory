package cardsupdater

import (
	"github.com/google/uuid"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
)

func convertProductResponseToCards(productsResponse ozon.ProductListInfoResponse) []card.Card {
	result := make([]card.Card, 0, len(productsResponse.Items))
	for _, item := range productsResponse.Items {
		img := ""
		if len(item.PrimaryImage) > 0 {
			img = item.PrimaryImage[0]
		}

		convertItem := card.Card{
			ID:          uuid.New(),
			Name:        item.Name,
			Article:     item.OfferId,
			Marketplace: card.MpOzon,
			IsComposite: false,
			Photo:       img,
		}

		result = append(result, convertItem)
	}

	return result
}
