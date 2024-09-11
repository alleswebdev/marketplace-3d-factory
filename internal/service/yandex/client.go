package yandex

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const (
	APIURL = "https://api.partner.market.yandex.ru"
)

type Client struct {
	httpClient http.Client
	token      string
	businessID string
	campaignID string
}

func NewClient(token, campaignID, businessID string) Client {
	return Client{
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
		token:      token,
		businessID: businessID,
		campaignID: campaignID,
	}
}

func (c Client) GetProductList(ctx context.Context) (OfferMappingsDTO, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", APIURL+"/businesses/"+c.businessID+"/offer-mappings", bytes.NewReader([]byte(`{"limit": 300, "archived": false}`)))
	if err != nil {
		return OfferMappingsDTO{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OfferMappingsDTO{}, errors.Wrap(err, "httpClient.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return OfferMappingsDTO{}, errors.Wrapf(err, "http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OfferMappingsDTO{}, errors.Wrap(err, "io.ReadAll")
	}

	var response OfferMappingsDTO
	err = json.Unmarshal(body, &response)
	if err != nil {
		return OfferMappingsDTO{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}

func (c Client) GetOrders(ctx context.Context, status string) (OrdersDTO, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", APIURL+"/campaigns/"+c.campaignID+"/orders?status="+status, nil)
	if err != nil {
		return OrdersDTO{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OrdersDTO{}, errors.Wrap(err, "httpClient.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return OrdersDTO{}, errors.Wrapf(err, "http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrdersDTO{}, errors.Wrap(err, "io.ReadAll")
	}

	var response OrdersDTO
	err = json.Unmarshal(body, &response)
	if err != nil {
		return OrdersDTO{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}
