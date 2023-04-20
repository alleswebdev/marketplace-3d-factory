package wb

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
	ApiURL = "https://suppliers-api.wildberries.ru"
)

type Client struct {
	httpClient http.Client
	token      string
}

func NewClient(token string) Client {
	return Client{
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

func (c Client) GetNewOrders(ctx context.Context) (OrdersResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", ApiURL+"/api/v3/orders/new", nil)
	if err != nil {
		return OrdersResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return OrdersResponse{}, errors.Wrap(err, "httpClient.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return OrdersResponse{}, errors.Wrapf(err, "http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrdersResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response OrdersResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return OrdersResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}

func (c Client) GetCardsList(ctx context.Context) (CardsListResponse, error) {
	requestBody, err := json.Marshal(CardListRequest{
		CardListSettings: CardListSettings{
			CardListCursor: CardListCursor{
				Limit: 1000,
			},
			Filter: Filter{
				WithPhoto: 1,
			},
		},
	})

	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "json.Marshal")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", ApiURL+"/content/v2/get/cards/list?locale=ru", bytes.NewReader(requestBody))
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "httpClient.Do")
	}

	if resp.StatusCode != http.StatusOK {
		return CardsListResponse{}, errors.Errorf("http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response CardsListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}
