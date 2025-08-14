package yandex

import (
	"context"
	"net/http"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/rest"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://api.partner.market.yandex.ru"
)

type Client struct {
	httpClient *rest.Client
	businessID string
	campaignID string
}

func NewClient(token, campaignID, businessID string) Client {
	return Client{
		httpClient: rest.NewClient(baseURL).WithBearerToken(token),
		businessID: businessID,
		campaignID: campaignID,
	}
}

func (c Client) GetProductList(ctx context.Context) (OfferMappingsDTO, error) {
	path := "/businesses/" + c.businessID + "/offer-mappings"
	resp, err := c.httpClient.DoRequest(ctx, http.MethodPost, path, GetProductListRequest{Limit: 300, Archived: false})
	if err != nil {
		return OfferMappingsDTO{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[OfferMappingsDTO](resp)
	if err != nil {
		return OfferMappingsDTO{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}

func (c Client) GetOrders(ctx context.Context, status string) (OrdersDTO, error) {
	path := "/campaigns/" + c.campaignID + "/orders?status=" + status
	resp, err := c.httpClient.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return OrdersDTO{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[OrdersDTO](resp)
	if err != nil {
		return OrdersDTO{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}
