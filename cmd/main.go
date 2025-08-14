package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alleswebdev/marketplace-3d-factory/internal/app/api"
	"github.com/alleswebdev/marketplace-3d-factory/internal/config"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/card"
	"github.com/alleswebdev/marketplace-3d-factory/internal/db/orderqueue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/ozon"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/queue"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/wb"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/cards_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/ozon_orders_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/supplies_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/wb_orders_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/workers/yandex_orders_updater"
	"github.com/alleswebdev/marketplace-3d-factory/internal/service/yandex"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.GetAppConfig()
	_ = os.Setenv("MARKETPLACE_APP_HOST", cfg.Host)

	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: false,
		StrictRouting: false,
		ServerHeader:  "go-app",
		AppName:       "Marketplace 3d factory",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1, http://localhost, http://127.0.0.1:4173, http://80.76.35.119",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Static("/", "./web/factory-front/dist")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wbClient := wb.NewClient(cfg.WbToken)
	ozonClient := ozon.NewClient(cfg.OzonToken, cfg.OzonClientID)
	yandexClient := yandex.NewClient(cfg.YandexToken, cfg.YandexCompaignID, cfg.YandexBusinessID)

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	cardStore := card.New(dbpool)
	orderQueueStore := orderqueue.New(dbpool)

	ordersUpdater := wb_orders_updater.NewWorker(wbClient, orderQueueStore, cardStore)
	go ordersUpdater.Run(ctx)

	ozonOrdersUpdater := ozon_orders_updater.NewWorker(ozonClient, orderQueueStore, cardStore)
	go ozonOrdersUpdater.Run(ctx)

	yandexOrdersUpdater := yandex_orders_updater.NewWorker(yandexClient, orderQueueStore, cardStore)
	go yandexOrdersUpdater.Run(ctx)

	suppliesUpdater := supplies_updater.NewWorker(wbClient, ozonClient, orderQueueStore)
	go suppliesUpdater.Run(ctx)

	cardsUpdater := cards_updater.NewWorker(wbClient, ozonClient, yandexClient, cardStore)
	go cardsUpdater.Run(ctx)

	queueService := queue.New(cardStore, orderQueueStore)
	appAPI := api.New(queueService)
	app.Get("/api/v2/list-queue", appAPI.ListQueue)
	app.Post("/api/v2/set-complete", appAPI.SetComplete)
	app.Post("/api/v2/set-children-complete", appAPI.SetChildrenComplete)
	app.Post("/api/v2/set-printing", appAPI.SetPrinting)

	err = app.Listen(":" + strconv.Itoa(cfg.Port))
	if err != nil {
		log.Fatal(err)
	}
}
