package api

import (
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
	resp, err := a.queueStore.GetAllItems(c.UserContext())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "queueStore.GetAllItems").Error())
	}

	return c.JSON(resp)
}
