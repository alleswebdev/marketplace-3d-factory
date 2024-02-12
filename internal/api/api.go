package api

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
)

type FactoryAPI struct {
	queueStore queue.Store
	cardStore  card.Store
	wbClient   wb.Client //todo перенести в воркер
}

func New(queueStore queue.Store, cardStore card.Store, wbClient wb.Client) FactoryAPI {
	return FactoryAPI{queueStore: queueStore, cardStore: cardStore, wbClient: wbClient}
}

type ListResponse struct {
	Items []QueueItem `json:"items"`
}

type QueueItem struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	Article     string           `json:"article"`
	Color       card.Color       `json:"color"`
	Size        card.Size        `json:"size"`
	Marketplace card.Marketplace `json:"marketplace"`
	Photo       string           `json:"photo"`
	IsPrinting  bool             `json:"is_printing"`
	IsComplete  bool             `json:"is_complete"`
	Children    []QueueItem      `json:"children"`
	TimePassed  string           `json:"time_passed"`
	IsComposite bool             `json:"is_composite"`
}

func (a FactoryAPI) ListQueue(c *fiber.Ctx) error {
	filter := queue.ListFilter{}
	err := c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "QueryParser").Error())
	}

	queueItems, err := a.queueStore.GetList(c.Context(), filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.GetAllItems").Error())
	}

	articles := make([]string, 0, len(queueItems))
	queueItemsByOrderIDMap := make(map[int64][]queue.Item)
	for _, item := range queueItems {
		articles = append(articles, item.Article)
		queueItemsByOrderIDMap[item.Parent] = append(queueItemsByOrderIDMap[item.Parent], item)
	}

	cards, err := a.cardStore.GetByArticlesMap(c.Context(), articles)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "cardStore.GetByArticlesMap").Error())
	}

	return c.JSON(ListResponse{Items: makeResponseItems(queueItems, cards, queueItemsByOrderIDMap)})
}

func makeResponseItems(items []queue.Item, cards map[string]card.Card, childrensMap map[int64][]queue.Item) []QueueItem {
	if len(items) <= 0 {
		return nil
	}

	result := make([]QueueItem, 0, len(items))
	for _, item := range items {
		if item.Parent != 0 {
			continue
		}

		childrens := makeItems(childrensMap[item.OrderID], cards)
		slices.SortFunc(childrens, func(a, b QueueItem) int {
			if a.Article > b.Article {
				return 1
			}
			return -1
		})
		card := cards[item.Article]

		result = append(result, QueueItem{
			ID:          item.ID,
			Name:        card.Name,
			Article:     item.Article,
			Color:       card.Color,
			Size:        card.Size,
			Marketplace: card.Marketplace,
			Photo:       card.Photo,
			IsPrinting:  item.IsPrinting,
			IsComplete:  item.IsComplete,
			TimePassed:  getTimeLeft(item.OrderCreatedAt),
			Children:    childrens,
			IsComposite: card.IsComposite,
		})
	}

	return result
}

func makeItems(items []queue.Item, cards map[string]card.Card) []QueueItem {
	if len(items) <= 0 {
		return nil
	}

	result := make([]QueueItem, 0, len(items))
	for _, item := range items {
		card := cards[item.Article]
		result = append(result, QueueItem{
			ID:          item.ID,
			Name:        card.Name,
			Article:     item.Article,
			Color:       card.Color,
			Size:        card.Size,
			Marketplace: card.Marketplace,
			Photo:       card.Photo,
			IsPrinting:  item.IsPrinting,
			IsComplete:  item.IsComplete,
			TimePassed:  getTimeLeft(item.OrderCreatedAt),
			IsComposite: card.IsComposite,
		})
	}

	return result
}

func getTimeLeft(orderCreatedAt time.Time) string {
	diff := time.Now().Sub(orderCreatedAt)
	hours := int(diff.Hours())
	return fmt.Sprintf("%d ч. %d мин.", hours, int(diff.Minutes())-hours*60)
}

type CompleteRequest struct {
	ID    int64 `json:"id"`
	State bool  `json:"state"`
}

func (a FactoryAPI) SetComplete(c *fiber.Ctx) error {
	req := CompleteRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "BodyParser").Error())
	}

	err = a.queueStore.SetComplete(c.Context(), req.ID, req.State)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetComplete").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) SetPrinting(c *fiber.Ctx) error {
	req := CompleteRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "BodyParser").Error())
	}

	err = a.queueStore.SetPrinting(c.Context(), req.ID, req.State)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetPrinting").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) UpdateCards(c *fiber.Ctx) error {
	cardsResp, err := a.wbClient.GetCardsList(c.Context())
	if err != nil {
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "wbClient.GetCardsList").Error())
		}
	}

	err = a.cardStore.AddCards(c.Context(), card.ConvertCards(cardsResp.Cards))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "cardStore.AddCards").Error())
	}

	return c.SendStatus(http.StatusOK)
}
