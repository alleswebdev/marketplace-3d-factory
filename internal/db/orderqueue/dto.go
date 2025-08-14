package orderqueue

import (
	"database/sql"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
)

type Order struct {
	ID             string       `db:"id"`
	Article        string       `db:"article"`
	Items          Items        `db:"order_composite_items"`
	Marketplace    string       `db:"marketplace"`
	OrderCreatedAt sql.NullTime `db:"order_created_at"`
	CreatedAt      sql.NullTime `db:"created_at"`
	UpdatedAt      sql.NullTime `db:"updated_at"`
	Info           Info         `db:"info"`
	IsComplete     bool         `db:"is_complete"`
	IsPrinting     bool         `db:"is_printing"`
}

type Info struct {
	OrderNumber     string    `json:"order_number"`
	OrderShipmentAt time.Time `json:"order_shipment_date"`
	Quantity        int32     `json:"quantity"`
}

type Item struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	IsComplete bool   `json:"is_complete"`
}

type Items []Item

type ListFilter struct {
	WithParentComplete   bool   `json:"withParentComplete"`
	WithChildrenComplete bool   `json:"withChildrenComplete"`
	Marketplace          string `json:"marketplace"`
}

func (f ListFilter) GetMarketplace() string {
	if len(f.Marketplace) == 0 {
		return card.MpWb.String()
	}

	return f.Marketplace
}
