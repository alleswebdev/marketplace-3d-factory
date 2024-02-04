package wb

import (
	"time"
)

type Order struct {
	Address               interface{} `json:"address"`
	DeliveryType          string      `json:"deliveryType"`
	SupplyID              string      `json:"supplyId"`
	User                  interface{} `json:"user"`
	OrderUID              string      `json:"orderUid"`
	Article               string      `json:"article"`
	Rid                   string      `json:"rid"`
	CreatedAt             time.Time   `json:"createdAt"`
	Offices               []string    `json:"offices"`
	Skus                  []string    `json:"skus"`
	ID                    int64       `json:"id"`
	WarehouseID           int64       `json:"warehouseId"`
	NmID                  int64       `json:"nmId"`
	ChrtID                int64       `json:"chrtId"`
	Price                 int64       `json:"price"`
	ConvertedPrice        int64       `json:"convertedPrice"`
	CurrencyCode          int64       `json:"currencyCode"`
	ConvertedCurrencyCode int64       `json:"convertedCurrencyCode"`
	CargoType             int64       `json:"cargoType"`
}

type OrdersResponse struct {
	Orders []Order `json:"orders"`
	Next   int     `json:"next"`
}

type Card struct {
	NmID        int    `json:"nmID"`
	ImtID       int    `json:"imtID"`
	NmUUID      string `json:"nmUUID"`
	SubjectID   int    `json:"subjectID"`
	SubjectName string `json:"subjectName"`
	VendorCode  string `json:"vendorCode"`
	Brand       string `json:"brand"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Photos      []struct {
		X288  string `json:"516x288"`
		Big   string `json:"big"`
		Small string `json:"small"`
	} `json:"photos"`
	Dimensions struct {
		Length int `json:"length"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"dimensions"`
	Characteristics []struct {
		ID    int         `json:"id"`
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"characteristics"`
	Sizes []struct {
		ChrtID   int      `json:"chrtID"`
		TechSize string   `json:"techSize"`
		WbSize   string   `json:"wbSize"`
		Skus     []string `json:"skus"`
	} `json:"sizes"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CardsListResponse struct {
	Cards                   []Card `json:"cards"`
	CardsListResponseCursor `json:"cursor"`
}

type CardsListResponseCursor struct {
	UpdatedAt time.Time `json:"updatedAt"`
	NmID      int       `json:"nmID"`
	Total     int       `json:"total"`
}

type CardListRequest struct {
	CardListSettings `json:"settings"`
}

type CardListSettings struct {
	CardListCursor `json:"cursor"`
	Filter         `json:"filter"`
}

type CardListCursor struct {
	Limit int `json:"limit"`
}

type Filter struct {
	WithPhoto int `json:"withPhoto"`
}
