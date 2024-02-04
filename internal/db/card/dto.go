package card

import (
	"database/sql"

	"github.com/google/uuid"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/sku"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

type Card struct {
	ID          uuid.UUID    `db:"id"`
	Name        string       `db:"name"`
	Article     string       `db:"article"`
	Articles    []string     `db:"articles"`
	Files       []string     `db:"files"`
	Color       Color        `db:"color"`
	Size        Size         `db:"size"`
	Marketplace Marketplace  `db:"marketplace"`
	IsComposite bool         `db:"is_composite"`
	Photo       string       `db:"photo"`
	CreatedAt   sql.NullTime `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type Color string
type Size string
type Marketplace string

const (
	ColorBlack  Color = "черный"
	ColorWhite  Color = "белый"
	ColorOrange Color = "оранжевый"
	ColorYellow Color = "желтый"
	ColorPink   Color = "розовый"
	ColorRed    Color = "красный"
	ColorPurple Color = "фиолетовый"
	ColorGreen  Color = "зелёный"
	ColorBlue   Color = "синий"

	SizeML       Size = "M"
	SizeL        Size = "L"
	SizeXL       Size = "XL"
	SizeStandart Size = "standart"

	MpWb   Marketplace = "wb"
	MpOzon Marketplace = "ozon"
)

func convertCards(wbCards []wb.Card) []Card {
	result := make([]Card, 0, len(wbCards))
	for _, item := range wbCards {
		convertItem := Card{
			ID:      uuid.MustParse(item.NmUUID),
			Name:    item.Title,
			Article: item.VendorCode,
		}

		if len(item.Photos) > 0 {
			convertItem.Photo = item.Photos[0].Big
		}

		result = append(result, convertItem)
	}

	return result
}

func convertCards2sku(wbCards []wb.Card) []sku.SKU {
	result := make([]sku.SKU, 0, len(wbCards))
	for _, item := range wbCards {
		convertItem := sku.SKU{
			NmID:        uuid.MustParse(item.NmUUID),
			Name:        item.Title,
			Articles:    []string{item.VendorCode},
			Color:       sku.ColorBlack,
			Size:        sku.SizeStandart,
			Marketplace: sku.MpWb,
			IsComposite: false,
		}

		result = append(result, convertItem)
	}

	return result
}
