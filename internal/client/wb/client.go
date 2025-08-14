package wb

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alleswebdev/marketplace-3d-factory/internal/client/rest"
	"github.com/pkg/errors"
)

const (
	marketplaceApiUrl = "https://marketplace-api.wildberries.ru"
	contentApiUrl     = "https://content-api.wildberries.ru"
	newOrdersPath     = "/api/v3/orders/new"
	supplyOrdersPath  = "/api/v3/supplies/%s/orders"
	ordersStatusPath  = "/api/v3/orders/status"
	getCardsPath      = "/content/v2/get/cards/list?locale=ru"
)

type Client struct {
	marketplaceClient *rest.Client
	contentClient     *rest.Client
}

func NewClient(token string) Client {
	return Client{
		marketplaceClient: rest.NewClient(marketplaceApiUrl).WithToken(token),
		contentClient:     rest.NewClient(contentApiUrl).WithToken(token),
	}
}

func (c Client) GetNewOrders(ctx context.Context) (OrdersResponse, error) {
	resp, err := c.marketplaceClient.DoRequest(ctx, http.MethodGet, newOrdersPath, nil)
	if err != nil {
		return OrdersResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[OrdersResponse](resp)
	if err != nil {
		return OrdersResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}

func (c Client) GetCardsList(ctx context.Context, cursor CardListCursor) (CardsListResponse, error) {
	resp, err := c.contentClient.DoRequest(ctx, http.MethodPost, getCardsPath, CardListRequest{
		CardListSettings: CardListSettings{
			CardListCursor: cursor,
			Filter: Filter{
				WithPhoto: 1,
			},
		},
	})
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[CardsListResponse](resp)
	if err != nil {
		return CardsListResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}

func (c Client) GetSupplies(ctx context.Context, next int) (SuppliesResponse, error) {
	resp, err := c.marketplaceClient.DoRequest(ctx, http.MethodGet, "/api/v3/supplies?limit=300&next="+strconv.Itoa(next), nil)
	if err != nil {
		return SuppliesResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[SuppliesResponse](resp)
	if err != nil {
		return SuppliesResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}

func (c Client) GetSupplyOrders(ctx context.Context, supply string) (SupplyOrdersResponse, error) {
	resp, err := c.marketplaceClient.DoRequest(ctx, http.MethodGet, fmt.Sprintf(supplyOrdersPath, supply), nil)
	if err != nil {
		return SupplyOrdersResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[SupplyOrdersResponse](resp)
	if err != nil {
		return SupplyOrdersResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}

func (c Client) GetOrdersStatus(ctx context.Context, orders []uint64) (OrderStatusResponse, error) {
	resp, err := c.marketplaceClient.DoRequest(ctx, http.MethodPost, ordersStatusPath, OrderStatusRequest{Orders: orders})
	if err != nil {
		return OrderStatusResponse{}, errors.Wrap(err, "doRequest")
	}
	defer resp.Body.Close()

	result, err := rest.ParseBody[OrderStatusResponse](resp)
	if err != nil {
		return OrderStatusResponse{}, errors.Wrap(err, "rest.ParseBody")
	}

	return result, nil
}
