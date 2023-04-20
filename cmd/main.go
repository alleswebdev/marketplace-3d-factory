package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/alleswebdev/marketplace-3d-factory/internal/command"
	"github.com/alleswebdev/marketplace-3d-factory/internal/config"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/order"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/savepoint"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/sku"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/orders_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/queuer"
)

func main() {
	cfg := config.GetAppConfig()
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "Fiber",
		AppName:       "Test App v1.0.1",
	})

	app.Static("/", "./web/commander-front/dist")
	app.Get("/api/list", func(c *fiber.Ctx) error {
		return c.JSON(cfg.Commands)
	})
	app.Post("/api/exec/:name", func(c *fiber.Ctx) error {
		name := c.Params("name")
		if len(name) == 0 {
			return c.SendStatus(500)
		}

		commandsMap := command.GetCommandsMapFromConfig(cfg)
		if cmd, ok := commandsMap[name]; ok {
			result, err := cmd.Start()
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, errors.Wrap(err, "cmd.Start").Error())
			}

			return c.JSON(result)
		}

		return c.SendStatus(fiber.StatusNotFound)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wbClient := wb.NewClient(cfg.WbToken)
	cardsResp, err := wbClient.GetCardsList(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	orderStore := order.New(dbpool)
	cardStore := card.New(dbpool)
	skusStore := sku.New(dbpool)
	queueStore := queue.New(dbpool)
	savepointsStore := savepoint.New(dbpool)

	if err != nil {
		log.Fatal(err)
	}

	err = cardStore.AddCards(ctx, convertCards(cardsResp.Cards))
	if err != nil {
		log.Fatal(err)
	}

	err = skusStore.SetNmIDs(ctx, convertCards2sku(cardsResp.Cards))
	if err != nil {
		log.Fatal(err)
	}

	ordersUpdater := orders_updater.NewWorker(wbClient, orderStore)
	go ordersUpdater.Run(ctx)
	fmt.Println("я какого то хуя тут")
	queueUpdater := queuer.NewWorker(dbpool, orderStore, savepointsStore, queueStore, cardStore)
	go queueUpdater.Run(ctx)

	err = app.Listen(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func convertCards(wbCards []wb.Card) []card.Card {
	result := make([]card.Card, 0, len(wbCards))
	for _, item := range wbCards {
		convertItem := card.Card{
			ID:      uuid.MustParse(item.NmUUID),
			Name:    item.Title,
			Article: item.VendorCode,
		}

		if len(item.Photos) > 0 {
			convertItem.Photo = item.Photos[0].Big
		}

		result = append(result, convertItem)
	}

	return result
}

func convertCards2sku(wbCards []wb.Card) []sku.SKU {
	result := make([]sku.SKU, 0, len(wbCards))
	for _, item := range wbCards {
		convertItem := sku.SKU{
			NmID:        uuid.MustParse(item.NmUUID),
			Name:        item.Title,
			Articles:    []string{item.VendorCode},
			Color:       sku.ColorBlack,
			Size:        sku.SizeStandart,
			Marketplace: sku.MpWb,
			IsComposite: false,
		}

		result = append(result, convertItem)
	}

	return result
}
