package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/domain"
	"github.com/alleswebdev/marketplace-3d-factory/internal/utils"
	"github.com/pkg/errors"
)

type (
	CardProvider interface {
		GetByArticlesMap(ctx context.Context, articles []string) (map[string]card.Card, error)
	}

	OrderProvider interface {
		GetOrders(ctx context.Context, filter orderqueue.ListFilter) ([]orderqueue.Order, error)
		SetComplete(ctx context.Context, id string, isComplete bool) error
		SetPrinting(ctx context.Context, id string, isPrinting bool) error
		SetChildrenComplete(ctx context.Context, id string, isComplete bool) error
	}
)

type (
	Queue struct {
		cardProvider  CardProvider
		orderProvider OrderProvider
	}
)

func New(cardProvider CardProvider, orderProvider OrderProvider) *Queue {
	return &Queue{
		cardProvider:  cardProvider,
		orderProvider: orderProvider,
	}
}

func (q Queue) SetComplete(ctx context.Context, id string, state bool) error {
	if err := q.orderProvider.SetComplete(ctx, id, state); err != nil {
		return errors.Wrap(err, "orderProvider.SetComplete")
	}

	return nil
}

func (q Queue) SetPrinting(ctx context.Context, id string, state bool) error {
	if err := q.orderProvider.SetPrinting(ctx, id, state); err != nil {
		return errors.Wrap(err, "orderProvider.SetPrinting")
	}

	return nil
}

func (q Queue) SetChildrenComplete(ctx context.Context, id string, state bool) error {
	if err := q.orderProvider.SetChildrenComplete(ctx, id, state); err != nil {
		return errors.Wrap(err, "orderProvider.SetChildrenComplete")
	}

	return nil
}

func (q Queue) ListQueue(ctx context.Context, withParent, withChildren bool, marketplace string) ([]domain.QueueItem, error) {
	orders, err := q.orderProvider.GetOrders(ctx, orderqueue.ListFilter{
		WithParentComplete:   withParent,
		WithChildrenComplete: withChildren,
		Marketplace:          marketplace,
	})

	if err != nil {
		return nil, errors.Wrap(err, "orderProvider.GetOrders")
	}

	articles := make([]string, 0, len(orders))
	for _, item := range orders {
		articles = append(articles, item.Article)
	}

	cards, err := q.cardProvider.GetByArticlesMap(ctx, articles)
	if err != nil {
		return nil, errors.Wrap(err, "cardProvider.GetByArticlesMap")
	}

	return makeItems(orders, cards), nil
}

func makeItems(orders []orderqueue.Order, cards map[string]card.Card) []domain.QueueItem {
	if len(orders) <= 0 {
		return nil
	}

	result := make([]domain.QueueItem, 0, len(orders))
	for _, order := range orders {
		currentCard := cards[order.Article]
		result = append(result, domain.QueueItem{
			ID:             order.ID,
			OrderID:        order.ID,
			Name:           currentCard.Name,
			Article:        order.Article,
			Marketplace:    currentCard.Marketplace,
			Photo:          currentCard.Photo,
			IsPrinting:     order.IsPrinting,
			IsComplete:     order.IsComplete,
			TimePassed:     getTimePassed(order.OrderCreatedAt.Time),
			ShipmentDate:   getShipmentDate(order.Info.OrderShipmentAt),
			IsComposite:    currentCard.IsComposite,
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
