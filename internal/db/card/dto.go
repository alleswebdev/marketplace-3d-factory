package card

import (
	"database/sql"

	"github.com/google/uuid"
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
