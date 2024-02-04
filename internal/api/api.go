package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
)

type FactoryApi struct {
	queueStore queue.Store
}

func New(queueStore queue.Store) FactoryApi {
	return FactoryApi{queueStore: queueStore}
}

func (a FactoryApi) ListQueue(c *fiber.Ctx) error {
	filter := queue.ListFilter{}
	err := c.QueryParser(&filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "QueryParser").Error())
	}

	resp, err := a.queueStore.GetList(c.Context(), filter)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.GetAllItems").Error())
	}

	return c.JSON(resp)
}

type CompleteRequest struct {
	ID    int64 `json:"id"`
	State bool  `json:"state"`
}

func (a FactoryApi) SetComplete(c *fiber.Ctx) error {
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

func (a FactoryApi) SetPrinting(c *fiber.Ctx) error {
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
