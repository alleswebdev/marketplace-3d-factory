package ozon

import (
	"time"
)

const (
	StatusAwaitingPackaging = "awaiting_packaging"
	StatusAwaitingDeliver   = "awaiting_deliver"
)

type ProductListResponse struct {
	Result struct {
		Items []struct {
			ProductID    int64  `json:"product_id"`
			OfferID      string `json:"offer_id"`
			IsFboVisible bool   `json:"is_fbo_visible"`
			IsFbsVisible bool   `json:"is_fbs_visible"`
			Archived     bool   `json:"archived"`
			IsDiscounted bool   `json:"is_discounted"`
		} `json:"items"`
		Total  int    `json:"total"`
		LastID string `json:"last_id"`
	} `json:"result"`
}

type ProductListInfoRequest struct {
	OfferID    []string `json:"offer_id"`
	ProductIDs []int64  `json:"product_id"`
	Sku        []int64  `json:"sku"`
}

type ProductListInfoResponse struct {
	Result struct {
		Items []struct {
			ID               int           `json:"id"`
			Name             string        `json:"name"`
			OfferID          string        `json:"offer_id"`
			Barcode          string        `json:"barcode"`
			BuyboxPrice      string        `json:"buybox_price"`
			CategoryID       int           `json:"category_id"`
			CreatedAt        time.Time     `json:"created_at"`
			Images           []string      `json:"images"`
			MarketingPrice   string        `json:"marketing_price"`
			MinOzonPrice     string        `json:"min_ozon_price"`
			OldPrice         string        `json:"old_price"`
			PremiumPrice     string        `json:"premium_price"`
			Price            string        `json:"price"`
			RecommendedPrice string        `json:"recommended_price"`
			MinPrice         string        `json:"min_price"`
			Sources          []interface{} `json:"sources"`
			Stocks           struct {
				Coming   int `json:"coming"`
				Present  int `json:"present"`
				Reserved int `json:"reserved"`
			} `json:"stocks"`
			Errors            []interface{} `json:"errors"`
			Vat               string        `json:"vat"`
			Visible           bool          `json:"visible"`
			VisibilityDetails struct {
				HasPrice      bool `json:"has_price"`
				HasStock      bool `json:"has_stock"`
				ActiveProduct bool `json:"active_product"`
				Reasons       struct {
				} `json:"reasons"`
			} `json:"visibility_details"`
			PriceIndex   string        `json:"price_index"`
			Images360    []interface{} `json:"images360"`
			ColorImage   string        `json:"color_image"`
			PrimaryImage string        `json:"primary_image"`
			Status       struct {
				State            string        `json:"state"`
				StateFailed      string        `json:"state_failed"`
				ModerateStatus   string        `json:"moderate_status"`
				DeclineReasons   []interface{} `json:"decline_reasons"`
				ValidationState  string        `json:"validation_state"`
				StateName        string        `json:"state_name"`
				StateDescription string        `json:"state_description"`
				IsFailed         bool          `json:"is_failed"`
				IsCreated        bool          `json:"is_created"`
				StateTooltip     string        `json:"state_tooltip"`
				ItemErrors       []interface{} `json:"item_errors"`
				StateUpdatedAt   time.Time     `json:"state_updated_at"`
			} `json:"status"`
			State            string `json:"state"`
			ServiceType      string `json:"service_type"`
			FboSku           int    `json:"fbo_sku"`
			FbsSku           int    `json:"fbs_sku"`
			CurrencyCode     string `json:"currency_code"`
			IsKgt            bool   `json:"is_kgt"`
			Rating           string `json:"rating"`
			DiscountedStocks struct {
				Coming   int `json:"coming"`
				Present  int `json:"present"`
				Reserved int `json:"reserved"`
			} `json:"discounted_stocks"`
			IsDiscounted      bool      `json:"is_discounted"`
			HasDiscountedItem bool      `json:"has_discounted_item"`
			Barcodes          []string  `json:"barcodes"`
			UpdatedAt         time.Time `json:"updated_at"`
			PriceIndexes      struct {
				PriceIndex        string `json:"price_index"`
				ExternalIndexData struct {
					MinimalPrice         string  `json:"minimal_price"`
					MinimalPriceCurrency string  `json:"minimal_price_currency"`
					PriceIndexValue      float64 `json:"price_index_value"`
				} `json:"external_index_data"`
				OzonIndexData struct {
					MinimalPrice         string `json:"minimal_price"`
					MinimalPriceCurrency string `json:"minimal_price_currency"`
					PriceIndexValue      int    `json:"price_index_value"`
				} `json:"ozon_index_data"`
				SelfMarketplacesIndexData struct {
					MinimalPrice         string `json:"minimal_price"`
					MinimalPriceCurrency string `json:"minimal_price_currency"`
					PriceIndexValue      int    `json:"price_index_value"`
				} `json:"self_marketplaces_index_data"`
			} `json:"price_indexes"`
			Sku                   int `json:"sku"`
			DescriptionCategoryID int `json:"description_category_id"`
			TypeID                int `json:"type_id"`
		} `json:"items"`
	} `json:"result"`
}

type UnfulfilledListRequest struct {
	Dir    string `json:"dir"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`

	Filter UnfulfilledListRequestFilter `json:"filter"`
	With   UnfulfilledListRequestWith   `json:"with"`
}

type UnfulfilledListRequestFilter struct {
	CutoffFrom       time.Time     `json:"cutoff_from"`
	CutoffTo         time.Time     `json:"cutoff_to"`
	DeliveryMethodID []interface{} `json:"delivery_method_id"`
	ProviderID       []interface{} `json:"provider_id"`
	Status           string        `json:"status"`
	WarehouseID      []interface{} `json:"warehouse_id"`
}

type UnfulfilledListRequestWith struct {
	AnalyticsData bool `json:"analytics_data"`
	Barcodes      bool `json:"barcodes"`
	FinancialData bool `json:"financial_data"`
	Translit      bool `json:"translit"`
}

type UnfulfilledListResponse struct {
	Result struct {
		Postings []struct {
			PostingNumber  string `json:"posting_number"`
			OrderID        int64  `json:"order_id"`
			OrderNumber    string `json:"order_number"`
			Status         string `json:"status"`
			DeliveryMethod struct {
				ID            int64  `json:"id"`
				Name          string `json:"name"`
				WarehouseID   int64  `json:"warehouse_id"`
				Warehouse     string `json:"warehouse"`
				TplProviderID int    `json:"tpl_provider_id"`
				TplProvider   string `json:"tpl_provider"`
			} `json:"delivery_method"`
			TrackingNumber     string      `json:"tracking_number"`
			TplIntegrationType string      `json:"tpl_integration_type"`
			InProcessAt        time.Time   `json:"in_process_at"`
			ShipmentDate       time.Time   `json:"shipment_date"`
			DeliveringDate     interface{} `json:"delivering_date"`
			Cancellation       struct {
				CancelReasonID           int    `json:"cancel_reason_id"`
				CancelReason             string `json:"cancel_reason"`
				CancellationType         string `json:"cancellation_type"`
				CancelledAfterShip       bool   `json:"cancelled_after_ship"`
				AffectCancellationRating bool   `json:"affect_cancellation_rating"`
				CancellationInitiator    string `json:"cancellation_initiator"`
			} `json:"cancellation"`
			Customer interface{} `json:"customer"`
			Products []struct {
				Price         string   `json:"price"`
				OfferID       string   `json:"offer_id"`
				Name          string   `json:"name"`
				Sku           int      `json:"sku"`
				Quantity      int      `json:"quantity"`
				MandatoryMark []string `json:"mandatory_mark"`
				CurrencyCode  string   `json:"currency_code"`
			} `json:"products"`
			Addressee interface{} `json:"addressee"`
			Barcodes  struct {
				UpperBarcode string `json:"upper_barcode"`
				LowerBarcode string `json:"lower_barcode"`
			} `json:"barcodes"`
			AnalyticsData struct {
				Region               string    `json:"region"`
				City                 string    `json:"city"`
				DeliveryType         string    `json:"delivery_type"`
				IsPremium            bool      `json:"is_premium"`
				PaymentTypeGroupName string    `json:"payment_type_group_name"`
				WarehouseID          int64     `json:"warehouse_id"`
				Warehouse            string    `json:"warehouse"`
				TplProviderID        int       `json:"tpl_provider_id"`
				TplProvider          string    `json:"tpl_provider"`
				DeliveryDateBegin    time.Time `json:"delivery_date_begin"`
				DeliveryDateEnd      time.Time `json:"delivery_date_end"`
				IsLegal              bool      `json:"is_legal"`
			} `json:"analytics_data"`
			FinancialData struct {
				Products []struct {
					CommissionAmount     int         `json:"commission_amount"`
					CommissionPercent    int         `json:"commission_percent"`
					Payout               int         `json:"payout"`
					ProductID            int         `json:"product_id"`
					OldPrice             int         `json:"old_price"`
					Price                int         `json:"price"`
					TotalDiscountValue   int         `json:"total_discount_value"`
					TotalDiscountPercent float64     `json:"total_discount_percent"`
					Actions              []string    `json:"actions"`
					Picking              interface{} `json:"picking"`
					Quantity             int         `json:"quantity"`
					ClientPrice          string      `json:"client_price"`
					ItemServices         struct {
						MarketplaceServiceItemFulfillment                int `json:"marketplace_service_item_fulfillment"`
						MarketplaceServiceItemPickup                     int `json:"marketplace_service_item_pickup"`
						MarketplaceServiceItemDropoffPvz                 int `json:"marketplace_service_item_dropoff_pvz"`
						MarketplaceServiceItemDropoffSc                  int `json:"marketplace_service_item_dropoff_sc"`
						MarketplaceServiceItemDropoffFf                  int `json:"marketplace_service_item_dropoff_ff"`
						MarketplaceServiceItemDirectFlowTrans            int `json:"marketplace_service_item_direct_flow_trans"`
						MarketplaceServiceItemReturnFlowTrans            int `json:"marketplace_service_item_return_flow_trans"`
						MarketplaceServiceItemDelivToCustomer            int `json:"marketplace_service_item_deliv_to_customer"`
						MarketplaceServiceItemReturnNotDelivToCustomer   int `json:"marketplace_service_item_return_not_deliv_to_customer"`
						MarketplaceServiceItemReturnPartGoodsCustomer    int `json:"marketplace_service_item_return_part_goods_customer"`
						MarketplaceServiceItemReturnAfterDelivToCustomer int `json:"marketplace_service_item_return_after_deliv_to_customer"`
					} `json:"item_services"`
					CurrencyCode string `json:"currency_code"`
				} `json:"products"`
				PostingServices struct {
					MarketplaceServiceItemFulfillment                int `json:"marketplace_service_item_fulfillment"`
					MarketplaceServiceItemPickup                     int `json:"marketplace_service_item_pickup"`
					MarketplaceServiceItemDropoffPvz                 int `json:"marketplace_service_item_dropoff_pvz"`
					MarketplaceServiceItemDropoffSc                  int `json:"marketplace_service_item_dropoff_sc"`
					MarketplaceServiceItemDropoffFf                  int `json:"marketplace_service_item_dropoff_ff"`
					MarketplaceServiceItemDirectFlowTrans            int `json:"marketplace_service_item_direct_flow_trans"`
					MarketplaceServiceItemReturnFlowTrans            int `json:"marketplace_service_item_return_flow_trans"`
					MarketplaceServiceItemDelivToCustomer            int `json:"marketplace_service_item_deliv_to_customer"`
					MarketplaceServiceItemReturnNotDelivToCustomer   int `json:"marketplace_service_item_return_not_deliv_to_customer"`
					MarketplaceServiceItemReturnPartGoodsCustomer    int `json:"marketplace_service_item_return_part_goods_customer"`
					MarketplaceServiceItemReturnAfterDelivToCustomer int `json:"marketplace_service_item_return_after_deliv_to_customer"`
				} `json:"posting_services"`
				ClusterFrom string `json:"cluster_from"`
				ClusterTo   string `json:"cluster_to"`
			} `json:"financial_data"`
			IsExpress    bool `json:"is_express"`
			Requirements struct {
				ProductsRequiringGtd           []interface{} `json:"products_requiring_gtd"`
				ProductsRequiringCountry       []interface{} `json:"products_requiring_country"`
				ProductsRequiringMandatoryMark []interface{} `json:"products_requiring_mandatory_mark"`
				ProductsRequiringRnpt          []interface{} `json:"products_requiring_rnpt"`
				ProductsRequiringJwUin         []interface{} `json:"products_requiring_jw_uin"`
			} `json:"requirements"`
			ParentPostingNumber string   `json:"parent_posting_number"`
			AvailableActions    []string `json:"available_actions"`
			MultiBoxQty         int      `json:"multi_box_qty"`
			IsMultibox          bool     `json:"is_multibox"`
			Substatus           string   `json:"substatus"`
			PrrOption           string   `json:"prr_option"`
		} `json:"postings"`
		Count int `json:"count"`
	} `json:"result"`
}
