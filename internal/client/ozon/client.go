package ozon

import (
	"context"
	"net/http"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/rest"
	"github.com/pkg/errors"
)

const (
	baseURL         = "https://api-seller.ozon.ru"
	listPath        = "/v3/product/list"
	infoListPath    = "/v3/product/info/list"
	postingListPath = "/v3/posting/fbs/unfulfilled/list"
)

type Client struct {
	*rest.Client
}

func NewClient(apiKey, clientID string) Client {
	return Client{
		Client: rest.NewClient(baseURL).WithClientID(clientID).WithApiKey(apiKey),
	}
}

func (c Client) GetProductList(ctx context.Context, lastID string, limit int) (ProductListResponse, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, listPath, ProductListRequest{Limit: limit, LastId: lastID})
	if err != nil {
		return ProductListResponse{}, errors.Wrap(err, "doRequest")
	}

	response, err := rest.ParseBody[ProductListResponse](resp)
	if err != nil {
		return ProductListResponse{}, errors.Wrap(err, "rest.ParseBody")
	}
	resp.Body.Close()

	return response, nil
}

func (c Client) GetProductInfoList(ctx context.Context, productIDs []int64) (ProductListInfoResponse, error) {
	resp, err := c.DoRequest(ctx, http.MethodPost, infoListPath, ProductListInfoRequest{ProductIDs: productIDs})
	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[ProductListInfoResponse](resp)
	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}

func (c Client) GetUnfulfilledList(ctx context.Context, status string) (UnfulfilledListResponse, error) {
	const monthDuration = time.Hour * 24 * 30
	resp, err := c.DoRequest(ctx, http.MethodPost, postingListPath, UnfulfilledListRequest{
		Dir:   "ASC",
		Limit: 1000,
		Filter: UnfulfilledListRequestFilter{
			CutoffFrom: time.Now().Add(-monthDuration),
			CutoffTo:   time.Now().Add(monthDuration),
			Status:     status,
		},
	})
	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[UnfulfilledListResponse](resp)
	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}
