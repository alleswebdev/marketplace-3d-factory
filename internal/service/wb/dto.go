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
	UpdatedAt string `json:"updatedAt"`
	NmID      int    `json:"nmID"`
	Total     int    `json:"total"`
}

type CardListRequest struct {
	CardListSettings `json:"settings"`
}

type CardListSettings struct {
	CardListCursor `json:"cursor"`
	Filter         `json:"filter"`
}

type CardListCursor struct {
	UpdatedAt string `json:"updatedAt,omitempty"`
	NmID      int    `json:"nmID,omitempty"`
	Limit     int    `json:"limit"`
}

type Filter struct {
	WithPhoto int `json:"withPhoto"`
}

type SuppliesResponse struct {
	Supplies []struct {
		ClosedAt  time.Time   `json:"closedAt"`
		ScanDt    interface{} `json:"scanDt"`
		ID        string      `json:"id"`
		Name      string      `json:"name"`
		CreatedAt time.Time   `json:"createdAt"`
		CargoType int         `json:"cargoType"`
		Done      bool        `json:"done"`
	} `json:"supplies"`
	Next int `json:"next"`
}

type SupplyOrdersResponse struct {
	Orders []struct {
		User                  interface{} `json:"user"`
		OrderUID              string      `json:"orderUid"`
		Article               string      `json:"article"`
		Rid                   string      `json:"rid"`
		CreatedAt             time.Time   `json:"createdAt"`
		Offices               []string    `json:"offices"`
		Skus                  []string    `json:"skus"`
		ID                    int64       `json:"id"`
		WarehouseID           int         `json:"warehouseId"`
		NmID                  int         `json:"nmId"`
		ChrtID                int         `json:"chrtId"`
		Price                 int         `json:"price"`
		ConvertedPrice        int         `json:"convertedPrice"`
		CurrencyCode          int         `json:"currencyCode"`
		ConvertedCurrencyCode int         `json:"convertedCurrencyCode"`
		CargoType             int         `json:"cargoType"`
	} `json:"orders"`
}
