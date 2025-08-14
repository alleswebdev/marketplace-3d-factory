package domain

import (
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
)

type (
	QueueItem struct {
		ID             string            `json:"id"`
		OrderID        string            `json:"order_id"`
		Name           string            `json:"name"`
		Article        string            `json:"article"`
		Marketplace    card.Marketplace  `json:"marketplace"`
		Photo          string            `json:"photo"`
		IsPrinting     bool              `json:"is_printing"`
		IsComplete     bool              `json:"is_complete"`
		Children       []QueueItem       `json:"children"`
		TimePassed     string            `json:"time_passed"`
		ShipmentDate   string            `json:"shipment_date"`
		IsComposite    bool              `json:"is_composite"`
		Info           orderqueue.Info   `json:"info"`
		CompositeItems []orderqueue.Item `json:"composite_items"`
	}
)
