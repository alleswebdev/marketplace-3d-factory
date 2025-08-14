package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Client struct {
	httpClient http.Client
	baseURL    string

	apiKey      string
	clientID    string
	token       string
	bearerToken string
}

func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *Client) WithApiKey(value string) *Client {
	c.apiKey = value
	return c
}

func (c *Client) WithClientID(value string) *Client {
	c.clientID = value
	return c
}

func (c *Client) WithToken(value string) *Client {
	c.token = value
	return c
}

func (c *Client) WithBearerToken(value string) *Client {
	c.bearerToken = value
	return c
}

func (c *Client) DoRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
	req, err := c.makeRequest(ctx, method, path, body)
	if err != nil {
		return nil, errors.Wrap(err, "makeRequest")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "httpClient.Do")
	}

	return resp, nil
}

func (c *Client) makeRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal")
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "http.NewRequestWithContext")
	}

	if len(c.clientID) > 0 {
		req.Header.Set("Client-Id", c.clientID)
	}

	if len(c.apiKey) > 0 {
		req.Header.Set("Api-Key", c.apiKey)
	}

	if len(c.token) > 0 {
		req.Header.Set("Authorization", c.token)
	}

	if len(c.bearerToken) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func ParseBody[T any](resp *http.Response) (T, error) {
	var response T

	if resp.StatusCode != http.StatusOK {
		return response, errors.Errorf("http status:%d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response, errors.Wrap(err, "io.ReadAll")
	}

	if jsonErr := json.Unmarshal(body, &response); jsonErr != nil {
		return response, errors.Wrap(jsonErr, "json.Unmarshal")
	}

	return response, nil
}
