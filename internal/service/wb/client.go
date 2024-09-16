package wb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	APIURL = "https://suppliers-api.wildberries.ru"
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
	req, err := http.NewRequestWithContext(ctx, "GET", APIURL+"/api/v3/orders/new", nil)
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

func (c Client) GetCardsList(ctx context.Context, cursor CardListCursor) (CardsListResponse, error) {
	requestBody, err := json.Marshal(CardListRequest{
		CardListSettings: CardListSettings{
			CardListCursor: cursor,
			Filter: Filter{
				WithPhoto: 1,
			},
		},
	})

	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "json.Marshal")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", APIURL+"/content/v2/get/cards/list?locale=ru", bytes.NewReader(requestBody))
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	body, err := c.doRequest(req)
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "doRequest")
	}

	var response CardsListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}

func (c Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "httpClient.Do")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return body, errors.Wrap(err, "io.ReadAll")
}

func (c Client) GetSupplies(ctx context.Context, next int) (SuppliesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", APIURL+"/api/v3/supplies?limit=300&next="+strconv.Itoa(next), nil)
	if err != nil {
		return SuppliesResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SuppliesResponse{}, errors.Wrap(err, "httpClient.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return SuppliesResponse{}, errors.Wrapf(err, "http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SuppliesResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response SuppliesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return SuppliesResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}

func (c Client) GetSupplyOrders(ctx context.Context, supply string) (SupplyOrdersResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v3/supplies/%s/orders", APIURL, supply), nil)
	if err != nil {
		return SupplyOrdersResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Authorization", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SupplyOrdersResponse{}, errors.Wrap(err, "httpClient.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return SupplyOrdersResponse{}, errors.Wrapf(err, "http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SupplyOrdersResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response SupplyOrdersResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return SupplyOrdersResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}
