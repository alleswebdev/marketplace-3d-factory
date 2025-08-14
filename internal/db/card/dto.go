package card

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

type Card struct {
	ID          uuid.UUID    `db:"id"`
	Name        string       `db:"name"`
	Article     string       `db:"article"`
	Articles    []string     `db:"articles"`
	Files       []string     `db:"files"`
	Marketplace Marketplace  `db:"marketplace"`
	IsComposite bool         `db:"is_composite"`
	Photo       string       `db:"photo"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type Marketplace string

const (
	MpWb     Marketplace = "wb"
	MpOzon   Marketplace = "ozon"
	MpYandex Marketplace = "yandex"
)

func (m Marketplace) String() string {
	return string(m)
}

// nolint
func ConvertCards(wbCards []wb.Card) []Card {
	result := make([]Card, 0, len(wbCards))
	for _, item := range wbCards {
		convertItem := Card{
			ID:          uuid.MustParse(item.NmUUID),
			Name:        item.Title,
			Article:     item.VendorCode,
			Marketplace: MpWb,
			IsComposite: false,
		}

		if len(item.Photos) > 0 {
			convertItem.Photo = item.Photos[0].Big
		}

		result = append(result, convertItem)
	}

	return result
}
