package ozon

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
	APIURL = "https://api-seller.ozon.ru"
)

type Client struct {
	httpClient http.Client
	token      string
	clientID   string
}

func NewClient(token, clientID string) Client {
	return Client{
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
		token:    token,
		clientID: clientID,
	}
}

func (c Client) GetProductList(ctx context.Context) (ProductListResponse, error) {
	//todo limit и пагинация хардкодится, тк карточек мало
	req, err := http.NewRequestWithContext(ctx, "POST", APIURL+"/v2/product/list", bytes.NewReader([]byte(`{"last_id": "","limit": 300}`)))
	if err != nil {
		return ProductListResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Api-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ProductListResponse{}, errors.Wrap(err, "httpClient.Do")
	}
	if resp.StatusCode != http.StatusOK {
		return ProductListResponse{}, errors.Wrapf(err, "http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProductListResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response ProductListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ProductListResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}

func (c Client) GetProductInfoList(ctx context.Context, productIDs []int64) (ProductListInfoResponse, error) {
	requestBody, err := json.Marshal(ProductListInfoRequest{ProductIDs: productIDs})

	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "json.Marshal")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", APIURL+"/v2/product/info/list", bytes.NewReader(requestBody))
	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Api-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "httpClient.Do")
	}

	if resp.StatusCode != http.StatusOK {
		return ProductListInfoResponse{}, errors.Errorf("http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response ProductListInfoResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return ProductListInfoResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}

const monthDuration = time.Hour * 24 * 30

func (c Client) GetUnfulfilledList(ctx context.Context, status string) (UnfulfilledListResponse, error) {
	requestBody, err := json.Marshal(UnfulfilledListRequest{
		Dir:   "ASC",
		Limit: 1000,
		Filter: UnfulfilledListRequestFilter{
			CutoffFrom: time.Now().Add(-monthDuration),
			CutoffTo:   time.Now().Add(monthDuration),
			Status:     status,
		},
	})

	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "json.Marshal")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", APIURL+"/v3/posting/fbs/unfulfilled/list", bytes.NewReader(requestBody))
	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	req.Header.Set("Client-Id", c.clientID)
	req.Header.Set("Api-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "httpClient.Do")
	}

	if resp.StatusCode != http.StatusOK {
		return UnfulfilledListResponse{}, errors.Errorf("http status:%d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "io.ReadAll")
	}

	var response UnfulfilledListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return UnfulfilledListResponse{}, errors.Wrap(err, "json.Unmarshal")
	}

	return response, nil
}
