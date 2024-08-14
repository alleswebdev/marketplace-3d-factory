package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order_queue"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
	"github.com/alleswebdev/marketplace-3d-factory/internal/utils"
)

type FactoryAPI struct {
	cardStore  card.Store
	orderStore order_queue.Store
	wbClient   wb.Client   //todo перенести в воркер обновление карточек
	ozonClient ozon.Client //todo перенести в воркер
}

func New(cardStore card.Store, wbClient wb.Client, ozonClient ozon.Client, orderStore order_queue.Store) FactoryAPI {
	return FactoryAPI{cardStore: cardStore, wbClient: wbClient, ozonClient: ozonClient, orderStore: orderStore}
}

type ListResponse struct {
	Items []QueueItem `json:"items"`
}

type QueueItem struct {
	ID             string             `json:"id"`
	OrderID        string             `json:"order_id"`
	Name           string             `json:"name"`
	Article        string             `json:"article"`
	Color          card.Color         `json:"color"`
	Size           card.Size          `json:"size"`
	Marketplace    card.Marketplace   `json:"marketplace"`
	Photo          string             `json:"photo"`
	IsPrinting     bool               `json:"is_printing"`
	IsComplete     bool               `json:"is_complete"`
	Children       []QueueItem        `json:"children"`
	TimePassed     string             `json:"time_passed"`
	ShipmentDate   string             `json:"shipment_date"`
	IsComposite    bool               `json:"is_composite"`
	Info           order_queue.Info   `json:"info"`
	CompositeItems []order_queue.Item `json:"composite_items"`
}

type CompleteRequest struct {
	ID    string `json:"id"`
	State bool   `json:"state"`
}

type ChildrenCompleteRequest struct {
	ID    string `json:"id"`
	State bool   `json:"state"`
}

func (a FactoryAPI) SetCompleteV2(c *fiber.Ctx) error {
	req := CompleteRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "BodyParser").Error())
	}

	err = a.orderStore.SetComplete(c.Context(), req.ID, req.State)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetComplete").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) SetPrintingV2(c *fiber.Ctx) error {
	req := CompleteRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "BodyParser").Error())
	}

	err = a.orderStore.SetPrinting(c.Context(), req.ID, req.State)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetPrinting").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) SetChildrenCompleteV2(c *fiber.Ctx) error {
	req := ChildrenCompleteRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "BodyParser").Error())
	}

	err = a.orderStore.SetChildrenComplete(c.Context(), req.ID, req.State)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetComplete").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) ListQueueV2(c *fiber.Ctx) error {
	filter := order_queue.ListFilter{}
	err := c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "QueryParser").Error())
	}

	orders, err := a.orderStore.GetOrders(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "orderStore.GetOrdersByFilter").Error())
	}

	articles := make([]string, 0, len(orders))
	for _, item := range orders {
		articles = append(articles, item.Article)
	}

	cards, err := a.cardStore.GetByArticlesMap(c.Context(), articles)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "cardStore.GetByArticlesMap").Error())
	}

	return c.JSON(ListResponse{Items: makeResponseItemsV2(orders, cards)})

}

func makeResponseItemsV2(orders []order_queue.Order, cards map[string]card.Card) []QueueItem {
	if len(orders) <= 0 {
		return nil
	}

	result := make([]QueueItem, 0, len(orders))
	for _, order := range orders {
		card := cards[order.Article]

		result = append(result, QueueItem{
			ID:             order.ID,
			OrderID:        order.ID,
			Name:           card.Name,
			Article:        order.Article,
			Color:          card.Color,
			Size:           card.Size,
			Marketplace:    card.Marketplace,
			Photo:          card.Photo,
			IsPrinting:     order.IsPrinting,
			IsComplete:     order.IsComplete,
			TimePassed:     getTimePassed(order.OrderCreatedAt.Time),
			ShipmentDate:   getShipmentDate(order.Info.OrderShipmentAt),
			IsComposite:    card.IsComposite,
			Info:           order.Info,
			CompositeItems: order.Items,
		})
	}

	return result
}

func getTimePassed(orderCreatedAt time.Time) string {
	diff := time.Since(orderCreatedAt)
	hours := int(diff.Hours())
	return fmt.Sprintf("%d ч. %d мин.", hours, int(diff.Minutes())-hours*60)
}

func getShipmentDate(shipmentAt time.Time) string {
	if shipmentAt.IsZero() {
		return ""
	}

	month := utils.DeclensionGenitiveMonth(int32(shipmentAt.Month()))
	return fmt.Sprintf("%d %s", shipmentAt.Day(), month)
}
