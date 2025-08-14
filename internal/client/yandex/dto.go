package yandex

import "time"

type OfferMappingsDTO struct {
	Status string `json:"status"`
	Result struct {
		Paging struct {
		} `json:"paging"`
		OfferMappings []struct {
			Offer struct {
				OfferId               string   `json:"offerId"`
				Name                  string   `json:"name"`
				Category              string   `json:"category"`
				Pictures              []string `json:"pictures"`
				Vendor                string   `json:"vendor"`
				Barcodes              []string `json:"barcodes"`
				Description           string   `json:"description"`
				ManufacturerCountries []string `json:"manufacturerCountries"`
				WeightDimensions      struct {
					Length int     `json:"length"`
					Width  float64 `json:"width"`
					Height float64 `json:"height"`
					Weight float64 `json:"weight"`
				} `json:"weightDimensions"`
				VendorCode   string   `json:"vendorCode,omitempty"`
				Tags         []string `json:"tags,omitempty"`
				BoxCount     int      `json:"boxCount,omitempty"`
				Type         string   `json:"type,omitempty"`
				Downloadable bool     `json:"downloadable,omitempty"`
				Adult        bool     `json:"adult,omitempty"`
				BasicPrice   struct {
					Value        float64   `json:"value"`
					CurrencyId   string    `json:"currencyId"`
					DiscountBase float64   `json:"discountBase"`
					UpdatedAt    time.Time `json:"updatedAt"`
				} `json:"basicPrice"`
				PurchasePrice struct {
					Value      float64   `json:"value"`
					CurrencyId string    `json:"currencyId"`
					UpdatedAt  time.Time `json:"updatedAt"`
				} `json:"purchasePrice,omitempty"`
				AdditionalExpenses struct {
					Value      float64   `json:"value"`
					CurrencyId string    `json:"currencyId"`
					UpdatedAt  time.Time `json:"updatedAt"`
				} `json:"additionalExpenses,omitempty"`
				CofinancePrice struct {
					Value      float64   `json:"value"`
					CurrencyId string    `json:"currencyId"`
					UpdatedAt  time.Time `json:"updatedAt"`
				} `json:"cofinancePrice,omitempty"`
				CardStatus string `json:"cardStatus"`
				Campaigns  []struct {
					CampaignId int    `json:"campaignId"`
					Status     string `json:"status"`
				} `json:"campaigns"`
				SellingPrograms []struct {
					SellingProgram string `json:"sellingProgram"`
					Status         string `json:"status"`
				} `json:"sellingPrograms"`
				Archived             bool   `json:"archived"`
				CustomsCommodityCode string `json:"customsCommodityCode,omitempty"`
				Params               []struct {
					Name  string `json:"name"`
					Value string `json:"value"`
				} `json:"params,omitempty"`
			} `json:"offer"`
			Mapping struct {
				MarketSku          int64  `json:"marketSku"`
				MarketSkuName      string `json:"marketSkuName"`
				MarketModelId      int    `json:"marketModelId"`
				MarketModelName    string `json:"marketModelName"`
				MarketCategoryId   int    `json:"marketCategoryId"`
				MarketCategoryName string `json:"marketCategoryName"`
			} `json:"mapping"`
		} `json:"offerMappings"`
	} `json:"result"`
}

type OrdersDTO struct {
	Pager struct {
		Total       int `json:"total"`
		From        int `json:"from"`
		To          int `json:"to"`
		CurrentPage int `json:"currentPage"`
		PagesCount  int `json:"pagesCount"`
		PageSize    int `json:"pageSize"`
	} `json:"pager"`
	Orders []struct {
		Id                            int     `json:"id"`
		Status                        string  `json:"status"`
		Substatus                     string  `json:"substatus"`
		CreationDate                  string  `json:"creationDate"`
		UpdatedAt                     string  `json:"updatedAt"`
		Currency                      string  `json:"currency"`
		ItemsTotal                    float64 `json:"itemsTotal"`
		DeliveryTotal                 float64 `json:"deliveryTotal"`
		BuyerItemsTotal               float64 `json:"buyerItemsTotal"`
		BuyerTotal                    float64 `json:"buyerTotal"`
		BuyerItemsTotalBeforeDiscount float64 `json:"buyerItemsTotalBeforeDiscount"`
		BuyerTotalBeforeDiscount      float64 `json:"buyerTotalBeforeDiscount"`
		PaymentType                   string  `json:"paymentType"`
		PaymentMethod                 string  `json:"paymentMethod"`
		Fake                          bool    `json:"fake"`
		Items                         []struct {
			Id                       int     `json:"id"`
			OfferId                  string  `json:"offerId"`
			OfferName                string  `json:"offerName"`
			Price                    float64 `json:"price"`
			BuyerPrice               float64 `json:"buyerPrice"`
			BuyerPriceBeforeDiscount float64 `json:"buyerPriceBeforeDiscount"`
			PriceBeforeDiscount      float64 `json:"priceBeforeDiscount"`
			Count                    int     `json:"count"`
			Vat                      string  `json:"vat"`
			ShopSku                  string  `json:"shopSku"`
			Promos                   []struct {
				Type    string  `json:"type"`
				Subsidy float64 `json:"subsidy"`
			} `json:"promos,omitempty"`
			Subsidies []struct {
				Type   string  `json:"type"`
				Amount float64 `json:"amount"`
			} `json:"subsidies,omitempty"`
		} `json:"items"`
		Subsidies []struct {
			Type   string  `json:"type"`
			Amount float64 `json:"amount"`
		} `json:"subsidies,omitempty"`
		Delivery struct {
			Type                string `json:"type"`
			ServiceName         string `json:"serviceName"`
			DeliveryPartnerType string `json:"deliveryPartnerType"`
			Dates               struct {
				FromDate string `json:"fromDate"`
				ToDate   string `json:"toDate"`
				FromTime string `json:"fromTime"`
				ToTime   string `json:"toTime"`
			} `json:"dates"`
			Region struct {
				Id     int    `json:"id"`
				Name   string `json:"name"`
				Type   string `json:"type"`
				Parent struct {
					Id     int    `json:"id"`
					Name   string `json:"name"`
					Type   string `json:"type"`
					Parent struct {
						Id     int    `json:"id"`
						Name   string `json:"name"`
						Type   string `json:"type"`
						Parent struct {
							Id     int    `json:"id"`
							Name   string `json:"name"`
							Type   string `json:"type"`
							Parent struct {
								Id     int    `json:"id"`
								Name   string `json:"name"`
								Type   string `json:"type"`
								Parent struct {
									Id   int    `json:"id"`
									Name string `json:"name"`
									Type string `json:"type"`
								} `json:"parent,omitempty"`
							} `json:"parent,omitempty"`
						} `json:"parent"`
					} `json:"parent"`
				} `json:"parent"`
			} `json:"region"`
			Address struct {
				Country    string `json:"country"`
				Postcode   string `json:"postcode"`
				City       string `json:"city"`
				Street     string `json:"street"`
				House      string `json:"house"`
				Entrance   string `json:"entrance,omitempty"`
				Entryphone string `json:"entryphone,omitempty"`
				Floor      string `json:"floor,omitempty"`
				Apartment  string `json:"apartment,omitempty"`
				Gps        struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
				} `json:"gps"`
			} `json:"address"`
			DeliveryServiceId int     `json:"deliveryServiceId"`
			LiftPrice         float64 `json:"liftPrice"`
			Shipments         []struct {
				Id           int    `json:"id"`
				ShipmentDate string `json:"shipmentDate"`
				Boxes        []struct {
					Id           int    `json:"id"`
					FulfilmentId string `json:"fulfilmentId"`
				} `json:"boxes"`
			} `json:"shipments"`
			OutletCode string `json:"outletCode,omitempty"`
		} `json:"delivery"`
		Buyer struct {
			Type string `json:"type"`
		} `json:"buyer"`
		TaxSystem       string `json:"taxSystem"`
		CancelRequested bool   `json:"cancelRequested"`
		Notes           string `json:"notes,omitempty"`
	} `json:"orders"`
	Paging struct {
	} `json:"paging"`
}

type GetProductListRequest struct {
	Limit    int  `json:"limit"`
	Archived bool `json:"archived"`
}
