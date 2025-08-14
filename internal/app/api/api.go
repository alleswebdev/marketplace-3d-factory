package api

import (
	"context"
	"net/http"

	"github.com/alleswebdev/marketplace-3d-factory/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

type QueueService interface {
	SetComplete(ctx context.Context, id string, state bool) error
	SetPrinting(ctx context.Context, id string, state bool) error
	SetChildrenComplete(ctx context.Context, id string, state bool) error
	ListQueue(ctx context.Context, withParent, withChildren bool, marketplace string) ([]domain.QueueItem, error)
}

type FactoryAPI struct {
	queueService QueueService
}

func New(queue QueueService) FactoryAPI {
	return FactoryAPI{queueService: queue}
}

type (
	CompleteRequest struct {
		ID    string `json:"id"`
		State bool   `json:"state"`
	}

	ChildrenCompleteRequest struct {
		ID    string `json:"id"`
		State bool   `json:"state"`
	}

	ListResponse struct {
		Items []domain.QueueItem `json:"items"`
	}

	ListRequest struct {
		WithParentComplete   bool   `json:"withParentComplete"`
		WithChildrenComplete bool   `json:"withChildrenComplete"`
		Marketplace          string `json:"marketplace"`
	}
)

func (a FactoryAPI) SetComplete(c *fiber.Ctx) error {
	req := new(CompleteRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "BodyParser").Error())
	}

	if err := a.queueService.SetComplete(c.Context(), req.ID, req.State); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetComplete").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) SetPrinting(c *fiber.Ctx) error {
	req := new(CompleteRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "BodyParser").Error())
	}

	if err := a.queueService.SetPrinting(c.Context(), req.ID, req.State); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetPrinting").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) SetChildrenComplete(c *fiber.Ctx) error {
	req := new(ChildrenCompleteRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "BodyParser").Error())
	}

	if err := a.queueService.SetChildrenComplete(c.Context(), req.ID, req.State); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.SetComplete").Error())
	}

	return c.SendStatus(http.StatusOK)
}

func (a FactoryAPI) ListQueue(c *fiber.Ctx) error {
	filter := new(ListRequest)
	if err := c.QueryParser(filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, errors.Wrap(err, "QueryParser").Error())
	}

	items, err := a.queueService.ListQueue(c.Context(), filter.WithParentComplete, filter.WithChildrenComplete, filter.Marketplace)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "ListQueue").Error())
	}

	return c.JSON(ListResponse{Items: items})

}
